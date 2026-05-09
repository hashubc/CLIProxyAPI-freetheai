import json
import unittest
from unittest import mock


class DiscordCheckinTests(unittest.TestCase):
    def test_send_checkin_posts_expected_payload(self):
        from scripts import discord_checkin

        opener = mock.Mock()
        response = mock.MagicMock()
        response.__enter__.return_value.read.return_value = json.dumps({"id": "m1"}).encode()
        opener.open.return_value = response

        result = discord_checkin.send_checkin("token-value", opener=opener)

        self.assertEqual(result, {"id": "m1"})
        request = opener.open.call_args.args[0]
        self.assertEqual(
            request.full_url,
            "https://discord.com/api/v9/channels/1473159205048553705/messages",
        )
        self.assertEqual(request.get_method(), "POST")
        self.assertEqual(request.headers["Authorization"], "token-value")
        self.assertEqual(request.headers["Content-type"], "application/json")
        self.assertEqual(json.loads(request.data.decode()), {"content": "/checkin", "tts": False})


if __name__ == "__main__":
    unittest.main()
