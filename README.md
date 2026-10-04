<div align="center">

<img width="100%" alt="header" src="https://capsule-render.vercel.app/api?type=waving&height=210&text=Virus%20Game%20Bot&fontAlign=50&fontAlignY=36&fontSize=60&desc=Auto%20Daily%20Spin%20%7C%20Telegram%20Miniapp%20Bot&descAlign=50&descAlignY=58"/>

<img alt="typing" src="https://readme-typing-svg.demolab.com?font=Inter&size=18&duration=3000&pause=650&center=true&vCenter=true&width=900&lines=Auto+Daily+Spin;Auto+Daily+Case+Opening;Multi-Account+Support;Go+Helper+Powered"/>

<p>
  <img alt="python" src="https://img.shields.io/badge/Python-3.12+-3776AB?logo=python&logoColor=white"/>
  <img alt="go" src="https://img.shields.io/badge/Go-Helper-00ADD8?logo=go&logoColor=white"/>
  <img alt="telegram" src="https://img.shields.io/badge/Telegram-Miniapp-26A5E4?logo=telegram&logoColor=white"/>
  <img alt="multi-account" src="https://img.shields.io/badge/Multi--Account-Supported-111111"/>
  <img alt="license" src="https://img.shields.io/badge/by-Yuurisandesu-111111"/>
</p>

<p>
  <b>Virus Game Bot</b> is a full automation bot for the Virus Game Telegram Miniapp.<br/>
  It handles authentication, daily spin, and daily case opening automatically for multiple accounts.<br/>
  Built and distributed by <b>Yuurisandesu</b>.
</p>

</div>

---

##️ Requirements

- Python `3.12+`
- Go (to build the helper binary)

---

## Quick Start

**1. Clone this repository**

```bash
git clone https://github.com/Yuurisan-N1/Virus-Miniapp.git
cd Virus-Miniapp
```

**2. Install Python dependencies**

```bash
pip install -r requirements.txt
```

**3. Build the Go helper**

```bash
cd gohelper
go build -o virushelper ./cmd/virushelper
cd ..
```

> On Windows the output binary will be `virushelper.exe`

**4. Fill `data.txt` with your initData** 

**5. Run the bot**

```bash
python bot.py
```

---

## Configuration

**`data.txt` — Account initData**

Fill `data.txt` with Telegram WebApp `initData` for each account, one per line:

```
user=%7B%22id%22...&hash=abc123
user=%7B%22id%22...&hash=def456
```

> initData can be obtained from browser DevTools when opening the Virus Game Miniapp on Telegram Web.

---

## Features

### Authentication
- Auto login using Telegram initData
- Referral code applied automatically on first login

### Account Info
- Displays balance, stars balance, next spin time, and next case time after login

### Daily Spin
- Auto run daily free spin
- Displays prize obtained and story reward if available

### Daily Case
- Auto open daily free case
- Displays prize obtained per opening

### Multi-Account
- Load unlimited accounts from `data.txt`
- Each account processed independently per run

---

##️ File Structure

```text
Virus-Game-Miniapp/
├── gohelper/
│   ├── cmd/
│   │   └── virushelper/
│   │       └── main.go
│   ├── internal/
│   │   └── api/
│   │       ├── client.go
│   │       └── models.go
│   └── go.mod
├── bot.py
├── data.txt
└── requirements.txt
```

---

## Log Format

The bot uses colored terminal logs:

| Color | Meaning |
|-------|---------|
| 🟢 Green | Success / Info |
| 🟡 Yellow | Warning / skip |
| 🔴 Red | Error / failed |

---

## Disclaimer

This tool is built for educational and testing purposes only. Use it wisely and at your own responsibility. The developer is not responsible for any account issues or bans.

---

<div align="center">
<img width="100%" alt="footer" src="https://capsule-render.vercel.app/api?type=waving&height=120&section=footer"/>
</div>