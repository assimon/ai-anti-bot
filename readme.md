# ai-anti-bot

Anti-spam robot for [Telegram](https://telegram.org/) groups

### English | [简体中文](wiki/readme_zh.md)

<p align="center">
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-blue" alt="license MIT"></a>
<a href="https://golang.org"><img src="https://img.shields.io/badge/Golang-1.22.3-red" alt="Go version 1.22.3"></a>
<a href="https://github.com/tucnak/telebot"><img src="https://img.shields.io/badge/Telebot Framework-v3-lightgrey" alt="telebot v3"></a>
</p>

`Telegram` is a world-renowned, very convenient and elegant anonymous communication tool.        
However, due to the anonymity of the software, a lot of spam promotion information often appears in the group. We have no way to identify whether there is spam information in the group at all times.      
Fortunately, `Telegram` provides us with a very powerful `Api`. Now we can use `artificial intelligence` to automatically help us detect user behavior.     

If you are a group administrator of `Telegram`, you can directly deploy this project privately.     

If you are a `developer`, you can use this project to familiarize yourself with the interactive development of `Go language` and `Telegram`, so that you can use `Api` to develop your own robot later!


## References：
- Telegram Api：[Telegram Api](https://core.telegram.org/bots/api)      
- Telebot：[Telebot](https://github.com/tucnak/telebot)


## How to use

### Docker Compose

```shell
# Clone the project
git clone https://github.com/assimon/ai-anti-bot.git && cd ai-anti-bot && mkdir data

# Set up your configuration
cp config.yml.example data/config.yml

# Fill in the configuration according to your needs
vi config.yml

# start up
docker compose up -d
```
It's very simple, right?😄

## How to configure
```yml
telegram:
  proxy:
  token: ""     # Fill in your robot token
  groups: [""]  # Fill in the group id where the robot needs to take effect
  owners: [""]  # Fill in the super administrator's telegram user ID
identification_model: "chatgpt"

# The following is the judgment strategy. For example, 
# the meaning of the following configuration item is: If you have joined the group for more than 3 days and have spoken more than 3 times or have been verified once, you do not need to verify again.
# This is to save your tokens😊
strategy:
  joined_time: 3
  number_of_speeches: 3
  verification_times: 1
  
chatgpt:
  proxy: ""   # OpenAI's proxy address, if necessary
  apikey: ""  # apikey
  model: "gpt-4o-mini"   # The detection model to be used. Please note that versions below gpt4 do not support image and file interaction.

# If your native language is not Chinese but other languages, 
# please use the translation to replace the following prompt with the language you want.
prompt:
   ...
```

### Other commands
```
/start       # Survival detection: feedback will be given if the robot service is running normally

# We can also use the following command to add our own advertising button to the robot, but please be sure to configure the telegram.owners option

/add_ad     # /add_ad #Add a new ad, format: ad title|jump link|expiration time (with hours, minutes, seconds)|weight (in descending order, the larger the value, the higher the weight), for example: /add_ad Hello|https://google.com|2099-01-01 00:00:00|100

/all_ad     # View All Ads Button

/del_ad     # Delete an ad button, for example: /del_ad 1 (delete the ad with id 1)
```

## Preview
![preview.png](wiki/preview.png)