## **요청**

```rust
use reqwest;
use serde_json::json;
use std::error::Error;

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let client = reqwest::Client::new();
    let url = "http://3.35.13.7:8081/synthesize";

    let payload = json!({
        "service": "tts-v2-korean-hanil",
        "speaker": 0,
        "message": "테스트입니다",
        "speed": 1,
        "noise": 0,
        "play_type": "D",
        "script": false
    });

    let response = client
        .post(url)
        .header("Content-Type", "application/json")
        .json(&payload)
        .send()
        .await?;

    let text = response.text().await?;
    println!("{}", text);

    Ok(())
}
```



## Download
```rust
use reqwest;
use std::error::Error;
use std::fs::File;
use std::io::Write;
use std::path::Path;

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let url = "http://3.35.13.7:8081/link/20251013/5dce13d0-84ab-4a06-8012-0cc6ba4a187c";
    let save_path = "output.wav";

    let response = reqwest::get(url).await?;
    response.raise_for_status()?;

    let mut file = File::create(Path::new(save_path))?;
    let content = response.bytes().await?;
    file.write_all(&content)?;

    Ok(())
}
```