# golang_cicd

## **create .env file in parent folder golang_cicd**
## Adding line to test CICD
```
dockerRegistryUserID="vikas93/"
dockerRepoName="go-cicd"
PORT="9090"
PersonalAccessToken="Value of PAT"
#docker hub
Username="docker username"
Password="Access token or password"
ServerAddress="https://index.docker.io/v1/"
```

## Need to save data of registerApp api in below format
```
{

	"app_name": "Payment",
	"service": [
        {
			"name": "payment_sanjeev",
			"repourl": "https://github.com/opstree/sanjeev.git",
			"dockerfilepath": "attendance/",
			"dockerfilename": "Dockerfile"

		},
		{
			"name": "payment_vikash",
			"repourl": "https://github.com/opstree/vikash.git",
			"dockerfilepath": "attendance/",
			"dockerfilename": "Dockerfile"
		}

	]
}
```