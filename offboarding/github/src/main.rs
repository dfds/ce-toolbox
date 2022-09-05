extern crate core;

use futures::FutureExt;
use octocrab::{Octocrab, Error};
use serde::{Serialize, Deserialize};

#[tokio::main]
async fn main() {
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
    println!(":: DEPLOY KEYS ::");
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
                    for key in val.data {
                        println!("Key: {}", key.title);
                        println!("Url: https://github.com/dfds/{}/settings/keys", val.repo_name);
                    }
                    println!("\n");
                }
            },
            Err(err) => {
                panic!(err);
            }
        }
    }

    println!(":: Repository users ::");
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
                    println!("Repo: {}", val.repo_name);
                    for repo_user in val.data {
                        println!("Username: {}", repo_user.login);
                    }
                    println!("\n");
                }
            },
            Err(err) => {
                panic!(err);
            }
        }
    }
}

async fn get_key(octo : Octocrab, repo_name : String) -> Result<OctoFutureResp<Vec<Key>>, Error> {
    let route = format!("/repos/dfds/{}/keys", repo_name);
    let resp : Result<Vec<Key>, Error> = octo.get(route, None::<&()>).await;
    return match resp {
        Ok(val) => {
            Ok(OctoFutureResp {
                repo_name: repo_name.clone(),
                data: val
            })
        },
        Err(err) => Err(err)
    }
}

async fn get_collaborators(octo : Octocrab, repo_name : String) -> Result<OctoFutureResp<Vec<RepoCollaborator>>, Error> {
    let route = format!("/repos/dfds/{}/collaborators?affiliation=direct", repo_name);
    let resp = octo.get(route, None::<&()>).await;
    return match resp {
        Ok(val) => {
            Ok(OctoFutureResp {
                repo_name: repo_name.clone(),
                data: val
            })
        },
        Err(err) => Err(err)
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
    pub read_only : bool
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
