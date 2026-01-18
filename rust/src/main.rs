use axum::{
    extract::Json,
    http::StatusCode,
    response::{IntoResponse, Response},
    routing::post,
    Router,
};
use reqwest;
use serde::{Deserialize, Serialize};

#[tokio::main]
async fn main() {
    let app = Router::new().route("/targetEndpointCheck", post(target_endpoint_check));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

#[derive(Deserialize)]
struct TargetEndpoint {
    url: String,
}

#[derive(Serialize)]
struct TargetEndpointResult {
    url: String,
    result: bool,
}

async fn target_endpoint_check(Json(target_endpoint): Json<TargetEndpoint>) -> Response {
    let url = target_endpoint.url;
    match reqwest::get(&url).await {
        Ok(response) => {
            let result = response.status() != reqwest::StatusCode::NOT_FOUND;
            let target_endpoint_result = TargetEndpointResult {
                url: url,
                result: result,
            };
            (StatusCode::OK, Json(target_endpoint_result)).into_response()
        }
        Err(_) => {
            let target_endpoint_result = TargetEndpointResult {
                url: url,
                result: false,
            };
            (StatusCode::OK, Json(target_endpoint_result)).into_response()
        }
    }
}
