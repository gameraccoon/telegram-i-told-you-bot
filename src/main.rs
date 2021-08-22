use serde::{Deserialize, Serialize};
use std::env;
use std::fs::File;
use std::io::Read;
use std::panic;
use std::path::Path;
use teloxide::prelude::*;

#[tokio::main]
async fn main() {
    let default_panic = std::panic::take_hook();
    panic::set_hook(Box::new(move |panic_info| {
        if let Some(s) = panic_info.payload().downcast_ref::<&str>() {
            if s.starts_with("to_user: ") {
                println!("{}", &s[9..]);
                return;
            }
        } else if let Some(s) = panic_info.payload().downcast_ref::<String>() {
            if s.starts_with("to_user: ") {
                println!("{}", &s[9..]);
                return;
            }
        }

        default_panic(panic_info);
    }));

    run().await;
}

#[derive(Serialize, Deserialize, Debug)]
struct Configuration {
    telegram_api_token: String,
}

fn get_config_example_string() -> String {
    let config = Configuration {
        telegram_api_token: String::from("telegramtoken:data"),
    };
    serde_json::to_string_pretty(&config).unwrap()
}

async fn run() {
    let config_file_path = "config.json";

    if !Path::new(config_file_path).is_file() {
        println!(
            "Please create config.json file.\nExample:\n{}",
            get_config_example_string()
        );
        return;
    }

    let mut config_file =
        File::open(config_file_path).expect("to_user: Can't read file config.json");

    let mut json_string = String::new();
    config_file.read_to_string(&mut json_string).unwrap();

    let config: Configuration = match serde_json::from_str(&json_string) {
        Ok(val) => val,
        Err(msg) => {
            println!("Error parsing config.json: {}", msg);
            return;
        }
    };

    teloxide::enable_logging!();

    let bot = Bot::new(config.telegram_api_token).auto_send();

    let bot_info = bot.get_me().await.unwrap();
    let account_name = bot_info.user.username.unwrap();
    println!("Authorized on account {}", account_name);

    teloxide::repl(bot, |message| async move {
        message.answer("test").await?;
        respond(())
    })
    .await;
}
