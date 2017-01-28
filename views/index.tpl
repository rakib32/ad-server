<!DOCTYPE html>

<html>

<head>
    <title>Ad Server</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
</head>

<body>

    <h1>Welcome to Ad Server</h1>
    <h3> Ad Server API's</h3>
    <div>
        <h3>Ad Request API(GET)</h3>
        <p>This api will handle the ad request and will serve the ads based on different ad targetings.</p>
        <div>Click the following link to get the ads:
            <a target="_blank" href="http://localhost:8080/api/v1/ads?os=android&version=6.0&token=1234&placement_key=1&limit=1"><b>http://localhost:8080/api/v1/ads?os=android&version=6.0&token=1234&placement_key=1&limit=1</b></a>
            <p>The parameters list are the following:</p>
            <ul>
                <li><b>limit (optional)</b>: the number of ads to be served on response to this request. Default value is 2.</li>
                <li><b>* os (required)</b>: operative system. Valid options: ios / android</li>
                <li><b>* version (required)</b>: operative system version, in a readable format (such as 4.2, 7.3, 6.0...)</li>
                <li><b>* token (required)</b>: the unique identifier for the user. Riggt now it's not implemented. It's a future improvement</li>
                <li><b>* placement_key</b>: (required) token identifier for the ad placement. As there is no user interface for this but you can find adspaces id in the adspace table.</li>

            </ul>
            The response for the ads request will be a JSON file containing a list with the offers and their assets, such as:
            <pre>
              [
                {
                  "Adid": 3,
                  "AdName": "Whats App",
                  "AppStoreURL": "https://www.whatsapp.com",
                  "ClickURL": "http://localhost:8080/api/v1/ad/click/3/65",
                  "ImpressionURL": "http://localhost:8080/api/v1/ad/impression/3/65",
                  "ActionText": "INSTALL",
                  "Images": {
                    "Icon": {
                      "Url": "",
                      "Height": "",
                      "Width": ""
                    },
                    "Banner": {
                      "Url": "http://www.freedigitalphotos.net/images/img/homepage/87357.jpg",
                      "Height": "300",
                      "Width": "300"
                    }
                  }
                }
              ]
            </pre>

            <ul>
                <li><b>Adid</b>: this is a numeric unique identifier of Offer. Will be used for further interactions like impression or click</li>
                <li><b>AdName(campaign name)</b>: the title of the offer. Should be displayed next to the icon of the app, as it's the most representative information about the campaigns.</li>
                <li><b>AppStoreURL</b>: this is the preview url. It will be a Google Play URL in case of Android apps, a iTunes URL in case of iOS, and a web URL in case it's another kind of offer.</li>
                <li><b>ClickURL</b>: this URL needs to be invoked when the user clicks in the ad. It's a post url </li>
                <li><b>ImpressionURL</b>: this URL needs to be invoked when ad is in inscreen. It's a post url</li>
                <li><b>ActionText</b>: It can have values such as "Download", "Install", "Try now", "Know more", etc.</li>
                <li><b>images</b>: an structure containing more images. It will also contain HTML code if there is any</li>
            </ul>
            <div>
                <h3>Track Impression API(POST)</h3>
                <p>This api will handle the impression request and it will also create the charge for that impression .</p>
                Example URL(POST: use rest console to post it): <span><b>http://localhost:8080/api/v1/ad/impression/{adid}/{deliveryid}</b></span>
                <p>The parameters list are the following:</p>
                <uL>
                    <li><b>adid</b>: Offer id to update the delivery summary./li>
                        <li><b>deliveryid</b> : Delivery of this ad.</li>
                </ul>
                The response for the ads request will be a JSON file containing a list with the offers and their assets, such as:
                <pre> 
                      {
                        "Message": "Impression has been tracked successfully."
                      }
                  </pre>
            </div>
            <div>
                <h3>Track Click API(POST)</h3>
                <p>This api will handle the click request and it will also create the charge for that click .</p>
                Example URL(POST: use rest console to post it): <span><b>http://localhost:8080/api/v1/ad/click/{adid}/{deliveryid}</b></span>
                <p>The parameters list are the following:</p>
                <uL>
                    <li><b>adid</b>: Offer id to update the delivery summary.</li>
                    <li><b>deliveryid</b> : Delivery of this ad.</li>
                </ul>
                <pre>
                  {
                    "Message": "Click has been tracked successfully."
                  }
                 </pre>
            </div>
        </div>

</body>

</html>