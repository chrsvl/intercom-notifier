# intercom-notifier

Notify a channel that a new Intercom conversation has been assigned to a specific team.

### Setup
Declare these 3 environment variables:

+ *CONVERSATION_BASE_URL*: the path to the Intercom conversation panel https://app.intercom.io/a/apps/YOUR_APP_ID/respond/inbox/AN_INTERCOM_TEAM
+ *SLACK_WEBHOOK_URL*: the Slack webhook url you have setup. intercom-notifier will push the notifications on that endpoint.
+ *TEAM_NAME*: the Intercom team name you want to monitor.
