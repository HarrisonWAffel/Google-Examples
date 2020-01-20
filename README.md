# Google Exercises 
###### Author: Harrison Affel

## Summary 
A collection of examples for working with various Google Cloud Platform API's. Currently, 
examples for Google TTS (Text to Speech),
STT (Speech to Text), and the Gmail API are provided. All examples are written in Golang, 
however in the future I may add additional languages.   



## Setting Up the GCP API's  
For these examples to work you will need to create a [Google Cloud Platform account](https://cloud.google.com/free/?utm_source=google&utm_medium=cpc&utm_campaign=na-US-all-en-dr-bkws-all-all-trial-p-dr-1008072&utm_content=text-ad-lpsitelinkPexp1-any-DEV_c-CRE_353294881642-ADGP_Hybrid+%7C+AW+SEM+%7C+BKWS+%7C+US+%7C+en+%7C+PHR+~+UX+Test+~+gcp-KWID_43700044772255395-kwd-224234454&utm_term=KW_gcp-ST_gcp&gclid=EAIaIQobChMI27HOuuGS5wIVjIbACh183wbkEAAYASABEgKv-_D_BwE).
There is a free trial for anyone reluctant to submit their payment information. After creating your account you will need to enable the following API's. 

1. [Google Text To Speech](https://cloud.google.com/text-to-speech/)
2. [Google Speech To Text](https://cloud.google.com/speech-to-text/)
3. [Gmail API](https://developers.google.com/gmail/api)

After enabling the API's you will need to create a service account by navigating to the IAM & Admin page, and then to the service accounts tab. 
When prompted, generate and download the JSON service account credential key. Rename this key to `creds.json` and place it inside the `credentials` folder.


Next, if you haven't done so already, download and install the [Google Cloud Platform CLI Tool](https://cloud.google.com/sdk/).
After installing the tool it will ask if you would like to authenticate, select yes. When prompted for an account, input the newly created service account email address.



If you encounter any permissions errors please check that the newly created service account has the required permissions
and that your gcp cli configuration has the correct account settings. (`gcloud config list` will show you your current configuration)
