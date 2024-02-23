#  StarLabs - Galxe KYC Grabber 


![Logo](https://i.postimg.cc/ZKbQdVgL/galxe.jpg)

## [SEE ENGLISH VERSION BELOW ](https://github.com/0xStarLabs/StarLabs-Galxe-KYC-Grabber?tab=readme-ov-file#english-version)👇

## 🔗 Links
[![Telegram channel](https://img.shields.io/endpoint?url=https://runkit.io/damiankrawczyk/telegram-badge/branches/master?url=https://t.me/StarLabsTech)](https://t.me/StarLabsTech)
[![Telegram chat](https://img.shields.io/endpoint?url=https://runkit.io/damiankrawczyk/telegram-badge/branches/master?url=https://t.me/StarLabsChat)](https://t.me/StarLabsChat)

🔔 CHANNEL: https://t.me/StarLabsTech

💬 CHAT: https://t.me/StarLabsChat

💰 DONATION EVM ADDRESS: 0x620ea8b01607efdf3c74994391f86523acf6f9e1


## 🤖 | Функционал:

🔬 Бот входит в аккаунт Galxe и достает ссылку для прохождения KYC.

🟢 Поддержка приватных ключей и мнемонических фраз.

🟢 Сохранение ссылок и приватных ключей в таблицу.

🟢 Обычные и мобильные прокси.

🟢 Многопоточность.


## 🚀 Installation
```
git clone https://github.com/0xStarLabs/StarLabs-Galxe-KYC-Grabber.git

cd StarLabs-Galxe-KYC-Grabber

go build

# Перед началом работы настройте необходимые модули в файлах config.yaml и /data

```

## ⚙️ Config

| Name | Description |
| --- | --- |
| max_tasks_retries | Максимальное количество попыток при выполнении задания |
| pause_between_tasks | Пауза между каждым действием |
| pause_between_accounts | Пауза между каждым аккаунтом |
| account_range | Диапазон аккаунтов для работы |
| mobile_proxy | Мобильные прокси |
| change_ip_pause | Пауза после смены айпи мобильных прокси |

## 🗂️ Data

Данные в папке data:

| Name | Description |
| --- | --- |
| wallets.txt | Содержит кошельки |
| proxies.txt | Содержит прокси в формате user:pass@ip:port |
| ip_change_links.txt | Содержит ссылки для смены айпи мобильных прокси |

## Дисклеймер
Автоматизация учетных записей пользователей Galxe, также известных как самостоятельные боты, является нарушением Условий обслуживания и правил сообщества Galxe и приведет к закрытию вашей учетной записи (аккаунтов). Рекомендуется осмотрительность. Я не буду нести ответственность за ваши действия. Прочтите об Условиях обслуживания Galxe и Правилах сообщества.

Это программное обеспечение было написано как доказательство концепции того, что учетные записи Galxe могут быть автоматизированы и могут выполнять действия, выходящие за рамки обычных пользователей Galxe, чтобы Galxe мог вносить изменения. Авторы  освобождаются от любой ответственности, которую может повлечь за собой ваше использование.

## ENGLISH VERSION:

## 🤖 | Features :

🔬 Bot logs into Galxe account and grab the link to pass KYC.

🟢 Support for private keys and mnemonic phrases.

🟢 Saves links and private keys to XLSX file.

🟢 Residential and mobile proxies.

🟢 Multithreading.

## 🚀 Installation
```
git clone https://github.com/0xStarLabs/StarLabs-Galxe-KYC-Grabber.git

cd StarLabs-Galxe-KYC-Grabber

go build

# Before you start, configure the required modules in config.yaml and /data files

```

## ⚙️ Config

| Name | Description |
| --- | --- |
| max_tasks_retries | Maximum number of attempts to complete a task |
| pause_between_tasks | pause between each action |
| pause_between_accounts | pause between each account |
| account_range | range of accounts to work |
| mobile_proxy | mobile proxies |
| change_ip_pause | pause after changing the ip of mobile proxies |

## 🗂️ Data

Data in the data folder:

| Name | Description |
| --- | --- |
| wallets.txt | Contains wallets |
| proxies.txt | Contains proxies in the format user:pass@ip:port |
| ip_change_links.txt | Contains links to change mobile proxy IPs |


## Disclamer
Automating Galxe user accounts, also known as autonomous bots, is a violation of Galxe's Terms of Service and Community Guidelines and will result in the termination of your account(s). Discretion is advised. I will not be held responsible for your actions. Read about Galxe's Terms of Service and Community Guidelines.

This software was written as a proof of concept that Galxe accounts can be automated and can perform actions beyond the normal Galxe users so that Galxe can make changes. The authors are released from any liability that your use may entail.
