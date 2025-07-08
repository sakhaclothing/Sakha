# Deplong Golang CI/CD to Google Cloud Platform

This is a simple Golang Model-Controller template using [Functions Framework for Go](https://github.com/GoogleCloudPlatform/functions-framework-go) and mongodb.com as the database host. It is compatible with Google Cloud Function CI/CD deployment.

Start here: Just [Fork this repo](https://github.com/gocroot/gcp/)

## MongoDB Preparation

The first thing to do is prepare a Mongo database using this template:

1. Sign up for mongodb.com and create one instance of Data Services of mongodb.
2. Go to Network Access menu > + ADD IP ADDRESS > ALLOW ACCESS FROM ANYWHERE  
   ![image](https://github.com/gocroot/gcp/assets/11188109/a16c5a73-ccdc-4425-8333-73c6fbf78e6d)  
3. Download [MongoDB Compass](https://www.mongodb.com/try/download/compass), connect with your mongo string URI from mongodb.com
4. Create database name iteung and collection reply  
   ![image](https://github.com/gocroot/alwaysdata/assets/11188109/23ccddb7-bf42-42e2-baac-3d69f3a919f8)  
5. Import [this json](https://whatsauth.my.id/webhook/iteung.reply.json) into reply collection.  
   ![image](https://github.com/gocroot/alwaysdata/assets/11188109/7a807d96-430f-4421-95fe-1c6a528ba428)  
   ![image](https://github.com/gocroot/alwaysdata/assets/11188109/fd785700-7347-4f4b-b3b9-34816fc7bc53)  
   ![image](https://github.com/gocroot/alwaysdata/assets/11188109/ef236b4d-f8f9-42c6-91ff-f6a7d83be4fc)  
6. Create a profile collection, and insert this JSON document with your 30-day token and WhatsApp number.  
   ![image](https://github.com/gocroot/alwaysdata/assets/11188109/5b7144c3-3cdb-472b-8ab3-41fe86dad9cb)  
   ![image](https://github.com/gocroot/alwaysdata/assets/11188109/829ae88a-be59-46f2-bddc-93482d0a4999)  

   ```json
   {
      "token": "v4.public.asdasfafdfsdfsdf",
      "phonenumber": "62881022526506",
      "secret": "secretkamuyangpanjangdanrumit089u08j32",
      "url": "https://asia-southeast2-awangga.cloudfunctions.net/logiccoffee/webhook/nomor/62881022526506",
      "urlapitext": "https://api.wa.my.id/api/v2/send/message/text",
      "urlapiimage": "https://api.wa.my.id/api/send/message/image",
      "urlapidoc": "https://api.wa.my.id/api/send/message/document",
      "urlqrlogin": "https://api.wa.my.id/api/whatsauth/request",
      "qrkeyword": "wh4t5auth0",
      "publickey": "0d6171e848ee9efe0eca37a10813d12ecc9930d6f9b11d7ea594cac48648f022",
      "botname": "lofe",
      "triggerword": "lofe",
      "telegramtoken": "",
      "telegramname": ""
   }
   ```

   ![image](https://github.com/gocroot/alwaysdata/assets/11188109/06330754-9167-4bf4-a214-5d75dab7c60a)  

## Folder Structure

This boilerplate has several folders with different functions, such as:

* .github: GitHub Action yml configuration.
* config: all apps configuration like database, API, token.
* controller: all of the endpoints functions
* model: all of the type structs used in this app
* helper: helper folder with a list of functions only called by others file
* route: all routes URL

## GCP Cloud Function CI/CD setup

To get an auth in Google Cloud, you can do the following:

1. Open Cloud Shell Terminal, type this command line per line:  
   
   ![image](https://github.com/gocroot/gcp/assets/11188109/14f8e9d7-f74c-4f74-ab9c-72731a3e5f13)  

   ```sh
   # Get a List of Project IDs in Your GCP Account
   gcloud projects list --format="value(projectId)"
   # Set Project ID Variable
   PROJECT_ID=yourprojectid
   # Create a service account
   gcloud iam service-accounts create "whatsauth" --project "${PROJECT_ID}"
   # Create JSON key for GOOGLE_CREDENTIALS variable in GitHub repo
   gcloud iam service-accounts keys create "key.json" --iam-account "whatsauth@${PROJECT_ID}.iam.gserviceaccount.com"
   # Read the key JSON file and copy the output, including the curl bracket, go to step 5.
   cat key.json
   # Authorize service account to act as admin in Cloud Run service
   gcloud projects add-iam-policy-binding ${PROJECT_ID} --member=serviceAccount:whatsauth@${PROJECT_ID}.iam.gserviceaccount.com --role=roles/run.admin
   # Authorize service account to delete artifact registry
   gcloud projects add-iam-policy-binding ${PROJECT_ID} --member=serviceAccount:whatsauth@${PROJECT_ID}.iam.gserviceaccount.com --role=roles/artifactregistry.admin
   # Authorize service account to deploy cloud function
   gcloud projects add-iam-policy-binding ${PROJECT_ID} --member=serviceAccount:whatsauth@${PROJECT_ID}.iam.gserviceaccount.com --role=roles/cloudfunctions.developer
   ```

3. Open Menu Cloud Build>settings, select the Service Account created by step 1, and enable Cloud Function Developer.  
   ![image](https://github.com/gocroot/gcp/assets/11188109/3ebc81b6-18b7-4d44-90b4-0abf67f82d66)  
   ![image](https://github.com/gocroot/gcp/assets/11188109/d2628542-99a6-44ce-ba78-798c249e0f22)  
5. Go to the GitHub repository; in the settings, menu>secrets>action, add GOOGLE_CREDENTIALS vars with the value from the key.json file.
6. Add other Vars into the secret>action menu:  

   ```sh
   MONGOSTRING=mongodb+srv://user:pass@gocroot.wedrfs.mongodb.net/
   WAQRKEYWORD=yourkeyword
   WEBHOOKURL=https://asia-southeast1-PROJECT_ID.cloudfunctions.net/gocroot/webhook/inbox
   WEBHOOKSECRET=yoursecret
   WAPHONENUMBER=62811111
   ```
7. Edit function name in main.yml (optional).


## Upgrade Apps

If you want to upgrade apps, please delete (go.mod) and (go.sum) files first, then type the command in your terminal or cmd :

```sh
go mod init gocroot
go mod tidy
```
