<p align="center">
  <img src="https://www.gstatic.com/android/keyboard/emojikitchen/20210521/u1fa84/u1fa84_u1f48c.png" width="100" />
</p>
<p align="center">
    <h1 align="center">postfix-to-cloudflare</h1>
</p>
<p align="center">
    <em>This project serves as a Postfix transport for parsing .eml (RFC822 compliant) emails and forwarding them via HTTP to a Cloudflare worker, leveraging the MailChannels free email sending API.</em>
</p>
<p align="center">
	<img src="https://img.shields.io/badge/GNU%20Bash-4EAA25.svg?style=flat&logo=GNU-Bash&logoColor=white" alt="GNU%20Bash">
	<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat&logo=Go&logoColor=white" alt="Go">
</p>

---

## Usage

This tool can be used in two primary modes:

1. **Parse-Only Mode**: This mode is useful for testing and debugging. It parses the `.eml` file but does not forward it to any endpoint. To use this mode, run the following command:
   ```
   postfix-to-cloudflare --parse-only < mail.eml
   ```
2. **Forwarding Mode**: In this mode, the tool parses the `.eml` file and forwards it to the specified endpoint using the provided token for authentication. This mode is used for actual operation. To use this mode, run the following command:
   ```
   postfix-to-cloudflare --endpoint YOUR_ENDPOINT --token YOUR_TOKEN < mail.eml
   ```
   Replace `YOUR_ENDPOINT` with the endpoint URL where the emails should be forwarded, and `YOUR_TOKEN` with the authentication token required by the Cloudflare worker.

## Getting Started

### Requirements

Ensure you have the following dependencies installed on your system:

- **Go**: `version >= 1.16` (if running from source)
- **Cloudflare Worker**: Ensure you have a Cloudflare account and have set up a worker by following the instructions in the [Cloudflare Worker Setup Guide](https://developers.cloudflare.com/workers/get-started/guide). The worker should be running the code from [this repository](https://github.com/Sh4yy/cloudflare-email) to handle incoming emails forwarded by this project.

> **Optional Requirement**: Ensure Postfix is configured to use this project as a transport. For guidance on configuring Postfix, refer to the [Postfix configuration guide](#configuring-postfix).

### Obtaining the latest release

To obtain the most recent version of `postfix-to-cloudflare`, navigate to the [Releases](https://github.com/zayigo/postfix-to-cloudflare/releases) page on GitHub and select the pre-compiled binary that matches your system
architecture.

### Building from source

1. Clone the repository:

```sh
git clone https://github.com/zayigo/postfix-to-cloudflare
```

2. Change to the project directory:

```sh
cd postfix-to-cloudflare
```

3. Build the app:

```sh
go build -o postfix-to-cloudflare main/main.go
```

### Running

Use the following command to run from source:

```sh
go run main/main.go
```

### Tests

Use the following command to run tests:

```sh
go test ./tests/...
```

## Limitations

Currently, the system does **not support handling attachments**, as they are not supported by the **MailChannels API**. This limitation is by design, with the primary focus on the core functionality of **email forwarding**. However, a potential future improvement could involve **uploading attachments to Cloudflare R2** and **embedding a link** to these attachments within the email body.

For examples of emails that this project can handle, please refer to the `tests/samples` directory.

## Configuring Postfix

This project is designed to seamlessly integrate as the primary transport mechanism for Postfix, in a manner completely transparent to other applications.

Assuming you have installed the project to `/usr/local/bin/postfix-to-cloudflare`, follow these steps to configure Postfix to utilize it as a transport mechanism:

1. **Edit the Postfix Master Configuration File**

   Open the `/etc/postfix/master.cf` file in your editor:

   ```sh
   nano /etc/postfix/master.cf
   ```

   Append the following lines to the end of the file. This configuration sets up a new service named `cloudflare` that Postfix will use to send emails:

   ```
   cloudflare   unix  -       n       n       -       -       pipe
     flags=FR user=nobody argv=/usr/local/bin/postfix-to-cloudflare --token YOUR_TOKEN --endpoint YOUR_ENDPOINT
   ```

2. **Modify the Main Configuration File**

   Next, edit the `/etc/postfix/main.cf` file to specify how emails should be routed:

   ```sh
   nano /etc/postfix/main.cf
   ```

   Append the following lines to the file. These settings are optional and allow you to rewrite the `MAIL_FROM` address and define the transport maps:

   ```
   # only if you want to rewrite the MAIL_FROM
   sender_canonical_maps = regexp:/etc/postfix/sender_canonical
   transport_maps = hash:/etc/postfix/transport
   ```

3. **Optional Step: Rewrite the MAIL_FROM Address**

   If you wish to rewrite the `MAIL_FROM` address for all emails sent through Postfix, perform the following step:

   Open the `/etc/postfix/sender_canonical` file in your editor:

   ```sh
   nano /etc/postfix/sender_canonical
   ```

   Add the following line to rewrite the sender address to `example@example.com` for all outgoing emails:

   ```
   /./ example@example.com
   ```

4. **Define the Transport Mechanism**

   Define how emails should be transported by editing the `/etc/postfix/transport` file:

   ```sh
   nano /etc/postfix/transport
   ```

   Add the following line to specify that all emails should be handled by the `cloudflare` service defined earlier:

   ```
   *    cloudflare
   ```

5. **Update Postfix's Transport Maps**

   After modifying the transport configuration, update Postfix's transport maps with the following command:

   ```sh
   postmap /etc/postfix/transport
   ```

6. **Restart Postfix**

   Finally, apply the changes by restarting Postfix:

   ```sh
   systemctl restart postfix
   ```

7. **Test Your Configuration**

   To verify that your Postfix configuration works as expected, you can send a test email from the command line:

   ```sh
   echo "Test email body" | mail -s "Test Email Subject" recipient@example.com
   ```

## License

This project is licensed under the GNU General Public License v3.0. For more details, see the [LICENSE](https://github.com/zayigo/postfix-to-cloudflare/blob/main/LICENSE) file.
