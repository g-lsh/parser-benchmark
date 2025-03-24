use lol_html::{element, HtmlRewriter, Settings};
use nix::sys::resource::{getrusage, UsageWho};
use std::time::{Duration};
use lol_html::html_content::ContentType;
use thiserror::Error;

fn get_cpu_time() -> Duration {
    let usage = getrusage(UsageWho::RUSAGE_SELF).unwrap();
    let user_time = Duration::new(usage.user_time().tv_sec() as u64, usage.user_time().tv_usec() as u32 * 1000);
    let system_time = Duration::new(usage.system_time().tv_sec() as u64, usage.system_time().tv_usec() as u32 * 1000);

    user_time + system_time
}


#[derive(Error, Debug)]
pub enum CustomError {
    #[error("Reqwest error: {0}")]
    Reqwest(#[from] reqwest::Error),
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),
    #[error("Rewriting error: {0}")]
    Rewriting(#[from] Box<dyn std::error::Error>),
}
async fn fetch_and_process_html(url: &str) -> Result<Duration, CustomError> {
    let resp = if url.starts_with("https://") {
        // Fetch HTML
        reqwest::get(url).await?.bytes().await?
    } else {
        // Read HTML from local file
        let data = std::fs::read(url)?;
        bytes::Bytes::from(data)
    };

    println!("Body size: {} bytes", resp.len());
    println!("Body size: {:.2} MB", resp.len() as f64 / 1024.0 / 1024.0);

    let start_cpu = get_cpu_time();

    // Process HTML with streaming parser
    let mut modified_html = Vec::new();
    {
        let mut rewriter = HtmlRewriter::new(
            Settings {
                element_content_handlers: vec![element!("title", |el| {
                    el.set_inner_content("Modified Title", ContentType::Text);
                    Ok(())
                })],
                ..Settings::default()
            },
            |chunk: &[u8]| modified_html.extend_from_slice(chunk),
        );

        rewriter.write(&resp).map_err(|e| CustomError::Rewriting(Box::new(e)))?;
        rewriter.end().map_err(|e| CustomError::Rewriting(Box::new(e)))?;
    }

    let elapsed_cpu = get_cpu_time() - start_cpu;

    Ok(elapsed_cpu)
}

#[tokio::main]
async fn main() {
    let args: Vec<String> = std::env::args().collect();
    if args.len() < 2 {
        eprintln!("Usage: cargo run <URL>");
        std::process::exit(1);
    }

    let url = &args[1];

    match fetch_and_process_html(url).await {
        Ok(cpu_time) => {
            println!("CPU Time: {:?}", cpu_time);
        }
        Err(e) => {
            eprintln!("Error: {}", e);
        }
    }
}