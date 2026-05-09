#!/usr/bin/env python3
"""Send a /checkin message to the configured Discord channel."""

import argparse
import json
import os
import sys
import urllib.error
import urllib.request


DEFAULT_GUILD_ID = "1461555807731585158"
DEFAULT_CHANNEL_ID = "1473159205048553705"
DEFAULT_MESSAGE = "/checkin"
DISCORD_API_BASE = "https://discord.com/api/v9"
TOKEN_ENV_VAR = "DISCORD_TOKEN"


def send_checkin(
    token,
    *,
    channel_id=DEFAULT_CHANNEL_ID,
    message=DEFAULT_MESSAGE,
    opener=None,
):
    """Post a check-in message and return the Discord JSON response."""
    if not token:
        raise ValueError(f"{TOKEN_ENV_VAR} is required")

    payload = json.dumps({"content": message, "tts": False}).encode("utf-8")
    request = urllib.request.Request(
        f"{DISCORD_API_BASE}/channels/{channel_id}/messages",
        data=payload,
        headers={
            "Authorization": token,
            "Content-Type": "application/json",
        },
        method="POST",
    )

    opener = opener or urllib.request.build_opener()
    with opener.open(request) as response:
        body = response.read()

    if not body:
        return {}
    return json.loads(body.decode("utf-8"))


def parse_args(argv):
    parser = argparse.ArgumentParser(description="Send /checkin to a Discord channel.")
    parser.add_argument(
        "--token",
        default=os.getenv(TOKEN_ENV_VAR),
        help=f"Discord user/bot token. Defaults to ${TOKEN_ENV_VAR}.",
    )
    parser.add_argument(
        "--guild-id",
        default=DEFAULT_GUILD_ID,
        help="Discord guild/server id. Kept for clarity; Discord sends by channel id.",
    )
    parser.add_argument(
        "--channel-id",
        default=DEFAULT_CHANNEL_ID,
        help="Discord channel id.",
    )
    parser.add_argument(
        "--message",
        default=DEFAULT_MESSAGE,
        help="Message content to send.",
    )
    return parser.parse_args(argv)


def main(argv=None):
    args = parse_args(argv or sys.argv[1:])
    try:
        result = send_checkin(args.token, channel_id=args.channel_id, message=args.message)
    except ValueError as errValue:
        print(f"error: {errValue}", file=sys.stderr)
        return 2
    except urllib.error.HTTPError as errHTTP:
        detail = errHTTP.read().decode("utf-8", errors="replace")
        print(f"discord api error: {errHTTP.code} {detail}", file=sys.stderr)
        return 1
    except urllib.error.URLError as errURL:
        print(f"request error: {errURL}", file=sys.stderr)
        return 1

    message_id = result.get("id", "")
    if message_id:
        print(f"checkin sent: {message_id}")
    else:
        print("checkin sent")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
