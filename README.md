# secret-config

Storing secrets like API keys or username/passwords in plain text is not a good idea - especially for public repos. For local development having secrets already checked in for a test system is very convenient especially for new developers.

We can fullfil both requirements when we encrypt the credentials and store the secrect in a password vault like bitwarden. 

The idea of this project is to give a developer a convenient way to access the stored secret in the vault.

* Dev sets env variable DEV_MODE
* Dev starts applicaton
* Dev is asked for the secret in console, in parallel the given bitwarden URL to the secret is opened
* Secret is copy pasted from the vault
* App will continue the startup with the decrypted config

In Go we can add tags(other languages like java use annotations) to struct elements. For secret-config you can add the tag 'secret_url' which defines where to get the secret to decrypt the value. 

```
type Config struct {
	Secret1       string `secret_url:"https://bitwarden.ofyour.organization/#/vault?itemId=184e5343-442d-4aa9-ba1d-b13c007fe2b8"`
	Secret2       string `secret_url:"https://bitwarden.ofyour.organization/#/vault?itemId=184e5343-442d-4aa9-ba1d-b13c007fe2b8"`
	AnotherSecret string `secret_url:"https://bitwarden.ofyour.organization/#/vault?itemId=184e5343-442d-4aa9-ba1d-b13c007fe2b8"`
	NonSecret     string
}
```

When loading the configuration you add the function ```ScanforUrlAndOpen```. 

```
func loadConfig() Config {
	c := Config{}
	.... // load your config, e.g. json or env file
	c=  urltag.ScanForUrlAndOpen(c, "DEV_MODE")
	return c
}
```

To enable opening the browser and decryption set the enviromnent var defined as the second parameter of ```ScanForUrlAndOpen```. In the above sample it's "DEV_MODE"

# Status
It's currently a working proof of concept