# MailAPI

MailAPI is an easy-to-use self-hosted REST API for sending emails.

## Features

- send plain text emails
- send html emails
- send combined html + text emails
- add attachments

## Configuration

Configuration is done through the environment or a configuration file. Environment takes precedence over config file.

```bash
API_KEY=somesecretapikey # required
SMTP_HOST=localhost      # optional
SMTP_PORT=25             # optional
SMTP_TLS=false           # optional
SMTP_USERNAME=user       # optional
SMTP_PASSWORD=s3cr3t     # optional
HOST=127.0.0.1:8000      # optional
```

To use a config file start the app with the config file argument: `mailapi config.yml`.

```yaml
apiKey: somesecretApiKey
smtp:
  host: localhost
  port: 25
  tls: false
  username: user
  password: s3cr3t
host: 127.0.0.1:8000
```

## API

There's only one endpoint: `POST /api/send`

```json
{
  "from": {
    "name": "optional name",
    "address": "required email address"
  },
  "to": {
    "name": "optional name",
    "address": "required email address"
  },
  "subject": "required subject",
  "contentType": "content type",
  "content": "content",
  "textContent": "text content",
  "htmlContent": "html content",
  "attachments": [
    {
      "contentType": "required content type",
      "name": "required name of the attachment",
      "filename": "required filename of the attachment",
      "data": "required bas64 encoded file content"
    }
  ]
}
```

`content` takes precedence over {text,html}Content.

- to send a simple text mail, only set `textContent`
- to send an html-only mail, only set `htmlContent`
- to send a combined, text and html mail, set `textContent` and `htmlContent`
- if you want full control over the content, set the `content` parameter

The API requires authentication. Use Bearer token authentication when calling the API.

```
Authorization: Beader <API_KEY>
```