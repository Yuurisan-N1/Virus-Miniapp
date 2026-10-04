import os
import sys
import subprocess
from datetime import datetime
import pyfiglet
from colorama import init, Fore, Style

init(autoreset=False)

RESET = Style.RESET_ALL
BOLD = Style.BRIGHT
GREEN = Fore.GREEN
YELLOW = Fore.YELLOW
RED = Fore.RED
CYAN = Fore.CYAN
MAGENTA = Fore.MAGENTA

REF_CODE = "VQSbHMVHguvMM799"

HELPER_NAME = "virushelper.exe" if os.name == "nt" else "virushelper"
GO_HELPER_PATH = os.path.join("gohelper", HELPER_NAME)


def sanitize_text(text):
    value = str(text)
    result = []
    for ch in value:
        if ch in "[]#-":
            continue
        result.append(ch)
    return "".join(result)


def LG(message):
    print(f"{GREEN}{BOLD}INFO {sanitize_text(message)}{RESET}")


def LW(message):
    print(f"{YELLOW}{BOLD}WARN {sanitize_text(message)}{RESET}")


def LE(message):
    print(f"{RED}{BOLD}ERROR {sanitize_text(message)}{RESET}")


def banner():
    os.system("cls" if os.name == "nt" else "clear")
    ascii_art = pyfiglet.figlet_format("Yuurisandesu", font="standard")
    print(CYAN + BOLD + ascii_art + RESET)
    print(MAGENTA + BOLD + "Welcome to Yuuri, Virus Game Bot" + RESET)
    LG("Ready to hack the world?")
    print(f"{YELLOW}{BOLD}Current time: {datetime.now().strftime('%d-%m-%Y %H:%M:%S')}{RESET}\n")


def set_title():
    sys.stdout.write("\x1b]2;Virus Game Bot by : 佐賀県産 (YUURI)\x1b\\")
    sys.stdout.flush()


def format_prize_name_for_log(name):
    clean = sanitize_text(str(name))
    parts = clean.split(" ", 1)
    first = parts[0] if parts else ""
    if first.replace(".", "", 1).isdigit():
        value = first
        label = parts[1] if len(parts) > 1 else ""
        return f"{RED}{BOLD}{value}{GREEN}{BOLD}{' ' + label if label else ''}"
    return f"{GREEN}{BOLD}{clean}"


def load_init_data(path):
    if not os.path.exists(path):
        LE("Init data file is missing")
        return []
    entries = []
    with open(path, "r", encoding="utf-8") as handle:
        for raw in handle:
            line = raw.strip()
            if line:
                entries.append(line)
    if not entries:
        LE("Init data file does not contain entries")
        return []
    LG(f"Init data entries loaded count {len(entries)} total")
    return entries


def parse_helper_output(line):
    text = line.strip()
    if not text:
        return None
    parts = text.split("|")
    if not parts:
        return None
    status = parts[0].strip().upper()
    if status == "ERR":
        message = "|".join(parts[1:]).strip()
        return {"ok": False, "error": message}
    data = {"ok": True}
    for segment in parts[1:]:
        segment = segment.strip()
        if not segment:
            continue
        if "=" in segment:
            key, value = segment.split("=", 1)
            key = key.strip().lower()
            value = value.strip()
            data[key] = value
    return data


def run_helper(action, params):
    if not os.path.exists(GO_HELPER_PATH):
        LE("Go helper binary is missing build helper before running bot")
        return None
    args = [GO_HELPER_PATH, "-action", action]
    for key, value in params.items():
        if value is None:
            continue
        args.append(f"-{key}")
        args.append(str(value))
    try:
        result = subprocess.run(args, capture_output=True, text=True, encoding="utf-8")
    except Exception as exc:
        LE(f"Helper invocation failed reason {exc}")
        return None
    stdout = (result.stdout or "").strip()
    if not stdout:
        LE("Helper response body is empty")
        return None
    first_line = stdout.splitlines()[0].strip()
    parsed = parse_helper_output(first_line)
    if not parsed:
        LE("Helper response format is not valid")
        return None
    if not parsed.get("ok"):
        err_text = parsed.get("error", "unknown error")
        LE(f"Helper response returned error {err_text}")
        return None
    return parsed


def get_token(init_data):
    LG("Token request started")
    result = run_helper("auth", {"init": init_data, "ref": REF_CODE})
    if not result:
        LE("Token request did not complete successfully")
        return None
    token = result.get("token")
    if not token:
        LE("Token value is missing in helper response")
        return None
    LG("Token request completed")
    return token


def get_account_info(token):
    LG("Account profile request started")
    result = run_helper("me", {"token": token})
    if not result:
        LE("Account profile request did not complete successfully")
        return None
    first_name = result.get("firstname") or ""
    balance = result.get("balance")
    stars_balance = result.get("stars")
    next_spin = result.get("nextspin")
    next_case = result.get("nextcase")
    LG(f"Profile loaded for name {first_name}")
    if balance is not None:
        LG(f"Current balance value {balance}")
    if stars_balance is not None:
        LG(f"Current stars balance value {stars_balance}")
    if next_spin:
        LG(f"Next free spin time {next_spin}")
    if next_case:
        LG(f"Next case free spin time {next_case}")
    return result


def run_daily_spin(token):
    LG("Daily spin request started")
    result = run_helper("spin", {"token": token})
    if not result:
        LE("Daily spin request did not complete successfully")
        return
    success_flag = result.get("success", "").lower()
    if success_flag not in ("true", "1", "yes"):
        LW("Daily spin result is not successful")
        return
    prize_name = result.get("prize") or "Unknown"
    formatted = format_prize_name_for_log(prize_name)
    line = f"{GREEN}{BOLD}RESULT Daily spin prize {formatted}{RESET}"
    print(line)
    story_reward = result.get("story")
    if story_reward:
        print(f"{GREEN}{BOLD}RESULT Story reward value {RED}{BOLD}{story_reward}{GREEN}{BOLD}{RESET}")


def run_daily_case(token, case_id):
    LG("Daily case opening request started")
    result = run_helper("case", {"token": token, "id": case_id})
    if not result:
        LE("Daily case request did not complete successfully")
        return
    success_flag = result.get("success", "").lower()
    if success_flag not in ("true", "1", "yes"):
        LW("Daily case result is not successful")
        return
    prize_name = result.get("prize") or "Unknown"
    formatted = format_prize_name_for_log(prize_name)
    line = f"{GREEN}{BOLD}RESULT Daily case prize {formatted}{RESET}"
    print(line)


def process_account(index_value, init_data, case_id):
    LG(f"Account {index_value} processing started")
    token = get_token(init_data)
    if not token:
        LE(f"Account {index_value} token sequence did not finish")
        print("")
        return
    LG(f"Authentication completed for account {index_value}")
    get_account_info(token)
    run_daily_spin(token)
    run_daily_case(token, case_id)
    print("")


def main():
    banner()
    set_title()
    init_path = "data.txt"
    entries = load_init_data(init_path)
    if not entries:
        return
    case_id = 2
    index_value = 1
    for init_data in entries:
        process_account(index_value, init_data, case_id)
        index_value += 1


if __name__ == "__main__":
    main()
