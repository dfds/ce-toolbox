extern crate core;

use std::collections::HashMap;
use futures::FutureExt;
use octocrab::{Octocrab, Error};
use octocrab::models::App;
use serde::{Serialize, Deserialize};
use std::fs::File;
use std::io::{Error as IoError, Write};

#[tokio::main]
async fn main() {
    let mut repos : HashMap<String, RepoContainer> = HashMap::new();
    let mut repos_with_permission_issues : HashMap<String, i8> = HashMap::new();
    let octo = Octocrab::builder()
        .personal_token(std::env::var("GITHUB_TOKEN").expect("No GITHUB_TOKEN environment variable was found"))
        .build().expect("Unable to build Octocrab client");

    let mut current_page = octo
        .orgs("dfds".to_owned())
        .list_repos()
        .per_page(100)
        .send()
        .await.unwrap();

    let mut prs = current_page.take_items();

    while let Ok(Some(mut new_page)) = octo.get_page(&current_page.next).await {
        prs.extend(new_page.take_items());
        current_page = new_page;
    }
    // Deploy keys
    let mut key_futures = Vec::new();
    for repo in &prs {
        let fut = get_key(octo.clone(), repo.name.clone());
        key_futures.push(fut.boxed());
    }
    let keys_results = futures::future::join_all(key_futures).await;

    for future_result in keys_results {
        match future_result {
            Ok(val) => {
                if val.data.len() > 0 {
                    if !repos.contains_key(&val.repo_name) {
                        repos.insert(val.repo_name.clone(), RepoContainer::new());
                    }
                    let mut repo_container = repos.get_mut(&val.repo_name).unwrap();
                    for key in val.data {
                        repo_container.keys.push(key);
                    }
                }
            },
            Err(err) => {
                let perm_issue = octo_error_handler(err.octo_error);
                if perm_issue.is_some() {
                    repos_with_permission_issues.insert(err.repo.clone(), 1);
                }
            }
        }
    }

    // Repo users
    let mut users_futures = Vec::new();
    for repo in &prs {
        let fut = get_collaborators(octo.clone(), repo.name.clone());
        users_futures.push(fut.boxed());
    }
    let users_results = futures::future::join_all(users_futures).await;

    for future_result in users_results {
        match future_result {
            Ok(val) => {
                if val.data.len() > 0 {
                    if !repos.contains_key(&val.repo_name) {
                        repos.insert(val.repo_name.clone(), RepoContainer::new());
                    }
                    let mut repo_container = repos.get_mut(&val.repo_name).unwrap();
                    for repo_user in val.data {
                        repo_container.repo_users.push(repo_user);
                    }
                }
            },
            Err(err) => {
                let perm_issue = octo_error_handler(err.octo_error);
                if perm_issue.is_some() {
                    repos_with_permission_issues.insert(err.repo.clone(), 1);
                }
            }
        }
    }


    // Print report
    if repos.len() > 0 {
        let mut users_file = File::create("report_users.csv").unwrap();
        let mut deploy_keys_file = File::create("report_keys.csv").unwrap();
        write!(users_file, "repository,username,url\n").unwrap();
        write!(deploy_keys_file, "repository,title,added_by,creator_url,created_at,last_used\n").unwrap();
        for (k, repo) in repos {
            if repo.is_empty() {
                continue;
            }
            println!("{}", k);

            if repo.keys.len() > 0 {
                println!("  Deploy keys:");
                for key in repo.keys {
                    println!("    Title        : {}", key.title);
                    println!("    Added by     : {}", key.added_by);
                    println!("    Creator URL  : https://github.com/{}", key.added_by);
                    println!("    Created at   : {}", key.created_at);
                    println!("    Last used    : {}", key.last_used);
                    print!("\n");
                    write!(deploy_keys_file, "{},\"{}\",{},https://github.com/{},{},{}\n", k, key.title, key.added_by, key.added_by, key.created_at, key.last_used).unwrap();
                }
            }

            if repo.repo_users.len() > 0 {
                println!("  Users:");
                for user in repo.repo_users {
                    println!("    Username: {}", user.login);
                    println!("    URL     : {}", user.html_url);
                    write!(users_file, "{},{},{}\n", k, user.login, user.html_url).unwrap();
                }
            }

        }
    }


    if repos_with_permission_issues.len() > 0 {
        let mut file = File::create("report_permission_issues.csv").unwrap();
        write!(file, "repository\n").unwrap();

        println!("Unable to check the following repositories due to permission issues:");
        for (k, i) in repos_with_permission_issues {
            println!(" - {} (https://github.com/{})", k, k);
            write!(file, "{}\n", k).unwrap();
        }

    }

}

fn octo_error_handler(err : Error) -> Option<()> {
    let mut panic_time = false;
    match &err {
        Error::Http { source, backtrace } => {
            match source.status() {
                Some(status_code) => {
                    if status_code.as_u16() == 404 || status_code.as_u16() == 403 {
                        //println!("Unable to query repo {} for {}. This is likely a permissions issue. Skipping.", source.url().unwrap().path(), category)
                        return Some(());
                    } else {
                        panic_time = true;
                    }
                },
                None => {
                    panic_time = true;
                }
            }
        }
        Error::Other { source, backtrace } => {},
        _ => {
            panic_time = true;
        }
    }

    if panic_time {
        println!("{:?}", err);
        panic!("Unable to proceed.");
    }

    None
}

async fn get_key(octo : Octocrab, repo_name : String) -> Result<OctoFutureResp<Vec<Key>>, AppError> {
    let route = format!("/repos/dfds/{}/keys", repo_name);
    let resp : Result<Vec<Key>, Error> = octo.get(route, None::<&()>).await;
    return match resp {
        Ok(val) => {
            Ok(OctoFutureResp {
                repo_name: repo_name.clone(),
                data: val
            })
        },
        Err(err) => Err(AppError {
            repo: repo_name.clone(),
            octo_error: err
        })
    }
}

async fn get_collaborators(octo : Octocrab, repo_name : String) -> Result<OctoFutureResp<Vec<RepoCollaborator>>, AppError> {
    let route = format!("/repos/dfds/{}/collaborators?affiliation=direct", repo_name);
    let resp = octo.get(route, None::<&()>).await;
    return match resp {
        Ok(val) => {
            Ok(OctoFutureResp {
                repo_name: repo_name.clone(),
                data: val
            })
        },
        Err(err) => Err(AppError {
            repo: repo_name.clone(),
            octo_error: err
        })
    }
}

struct RepoContainer {
    keys : Vec<Key>,
    repo_users : Vec<RepoCollaborator>
}

impl RepoContainer {
    pub fn new() -> Self {
        Self {
            keys: Vec::new(),
            repo_users: Vec::new()
        }
    }
}

struct AppError {
    repo : String,
    octo_error : octocrab::Error
}

impl RepoContainer {
    pub fn is_empty(&self) -> bool {
        let mut empty = true;

        if self.keys.len() > 0 {
            empty = false;
        }

        if self.repo_users.len() > 0 {
            empty = false;
        }

        empty
    }
}

pub struct OctoFutureResp<T> {
    data : T,
    repo_name : String
}

#[derive(Serialize, Deserialize)]
pub struct Key {
    pub id : i64,
    pub key : String,
    pub url : String,
    pub title : String,
    pub verified : bool,
    pub created_at : String,
    pub last_used : String,
    pub read_only : bool,
    pub added_by : String
}

#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct RepoCollaborator {
    pub login: String,
    pub id: i64,
    #[serde(rename = "node_id")]
    pub node_id: String,
    #[serde(rename = "avatar_url")]
    pub avatar_url: String,
    #[serde(rename = "gravatar_id")]
    pub gravatar_id: String,
    pub url: String,
    #[serde(rename = "html_url")]
    pub html_url: String,
    #[serde(rename = "followers_url")]
    pub followers_url: String,
    #[serde(rename = "following_url")]
    pub following_url: String,
    #[serde(rename = "gists_url")]
    pub gists_url: String,
    #[serde(rename = "starred_url")]
    pub starred_url: String,
    #[serde(rename = "subscriptions_url")]
    pub subscriptions_url: String,
    #[serde(rename = "organizations_url")]
    pub organizations_url: String,
    #[serde(rename = "repos_url")]
    pub repos_url: String,
    #[serde(rename = "events_url")]
    pub events_url: String,
    #[serde(rename = "received_events_url")]
    pub received_events_url: String,
    #[serde(rename = "type")]
    pub type_field: String,
    #[serde(rename = "site_admin")]
    pub site_admin: bool,
    pub permissions: Permissions,
}

#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Permissions {
    pub admin: bool,
    pub push: bool,
    pub pull: bool,
}
