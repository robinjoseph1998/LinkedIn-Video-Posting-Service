
# FOR CHECKING ACCESS-TOKEN SCOPES AND DETAILS

curl --location --request POST 'https://www.linkedin.com/oauth/v2/introspectToken' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=<Application Client ID>' \
--data-urlencode 'client_secret=<Application Client Secret>' \
--data-urlencode 'token=<Token Value>'


# ************ Video Uploading Processes **************** #

# Step:1 Initializing Video

curl --location --request POST 'https://api.linkedin.com/rest/videos?action=initializeUpload'
--header 'LinkedIn-Version: {version number in the format YYYYMM}' \
--header 'X-RestLi-Protocol-Version: 2.0.0' \
-H 'Authorization: Bearer {INSERT_TOKEN}' \
--header 'Content-Type: application/json' \
--data-raw '{ "initializeUploadRequest": {
       "owner": "urn:li:organization:2414183",
       "fileSizeBytes": 1055736 ,
       "uploadCaptions": false,
       "uploadThumbnail": false
    }
}'

# Step:2 Upload Video

curl -v \
 -H "Content-Type:application/octet-stream" \
 --upload-file ~/Downloads/sample.mp4 \
"https://www.linkedin.com/dms-uploads/C5505AQH-oV1qvnFtKA/uploadedVideo?sau=aHR0cHM6Ly93d3cubGlua2VkaW4tZWkuY29tL2FtYnJ5L2FtYnJ5LXZpZGVvZWkvP3gtbGktYW1icnktZXA9QVFHVkdRS0FtS05oM2dBQUFYd19ObG1uZzVYcllXajEzZjIybXh4LW55SGVBclVKcE8y"

# Step:3 Finalize Video

curl --location --request POST 'https://api.linkedin.com/rest/videos?action=finalizeUpload' \
--header 'LinkedIn-Version: {version number in the format YYYYMM}' \
--header 'X-RestLi-Protocol-Version: 2.0.0' \
-H 'Authorization: Bearer {INSERT_TOKEN}' \
--header 'Content-Type: application/json' \
--data-raw '{"finalizeUploadRequest": {"video": "urn:li:video:C4E10AQEfKKMV9a1d-g", "uploadToken": "", "uploadedPartIds": ["/ambry-video/signedId/AQJxWLYHLkIhgQAAAYBHKrkfFtLeuylYvIQ4PIvA8y9vsNxOL81dFKPzDoelwAh3u1LtNoV1E0QgYbnJSzQYufWwcuDfeQFpp99kXzG3pmRCi3kDmUGRaP3VAqt1Gu527c5Av5A2blyHvduKarZJ288qeqO_JPO20ivmoEgBMr1LqzFnZIDvCvbJjNdtk_c3iH6NBv527jU18_cT5CWTds86xvXzUQurrlRvBmG71szPI9B8VarjT_XNIfmz6kfcego3Np3DR9s.bin", "/ambry-video/signedId/AQImHSPUJvN5iQAAAYBHKrummJ26JZIelWw8dSfka6zbEpviF6Ibu2NMaiSgI8B_VOieCVIHePqnLwWCFkZO3fUYG5GT_RLBgDZqP_Al0Dm1qzov4xvQiQPdgwRIm6J7uoN2LUpyAgJfmkfdl7e2_HdJhuA4xp2mOu5fa3icuQWHTU09_ZlDZfS014Mcla7n6sB-X0LkiyJlJa9Ku8plE12-7hvU0KLd2GMzTlIxOrl9t6uOsMG5d4kLzV5Jje0Aifey4UM0we0.bin", "/ambry-video/signedId/AQLN2Hlak3Kk2AAAAYBHKsKYlebOqPNKgpg5gF38dkxo79E8WsGYpJbbxI9VO48o67pj05ajTbqqT7zm_eFv5HzMc3cyvdeVP9obIEfbR2M2Kjx5zisV4rTzKqDPIXu3CpxjuwnTdxS3I5-Pl8BvO4hCFqtxqpO3LW81rhys8kvykvOdRWCazP3H_3M7cbEmXmHz1aNSrFR3SIPFBUIUDW4P4O7IrOlYILG9PqqcdMweWvIDBPH1TyVJx-Lxd6Pm0kvLGg7QOmA.bin", "/ambry-video/signedId/AQJ53DRVN7kWEwAAAYBHKsYPzzQojALApZNqdLjenwy7yeH1RKGauKJ4Y0WtkJViorMsjQCzeMpG2Q1JNZSa_YeLXvr9hrY3TR5gw2cZfMTpnQ0wikzRI43OF-9X0oVOuaNBWfGj729ROuh6QRudsNOHwx497FPQwD9t4O7KNlt-9xNk2TkAeHXv9GKOUqXkjIOcG-VODb_5KXUiXihi921cqdeEbhKiFjAa5XQKthQBAxBewVLq0d619vyyjva38IAqNBiF0WM.bin"]
}
}'

# Step: 4 Posting the Video

curl -X POST "https://api.linkedin.com/v2/ugcPosts" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
-H "X-Restli-Protocol-Version: 2.0.0" \
-H "Content-Type: application/json" \
-d '{
    "author": "urn:li:person:YOUR_AUTHOR_URN",
    "commentary": "Check out this amazing video!",
    "visibility": "PUBLIC",
    "distribution": {
        "feedDistribution": "MAIN_FEED",
        "targetEntities": [],
        "thirdPartyDistributionChannels": []
    },
    "content": {
        "media": {
            "title": "Amazing Video Title",
            "id": "urn:li:video:YOUR_VIDEO_URN"
        }
    },
    "lifecycleState": "PUBLISHED",
    "isReshareDisabledByAuthor": false
}'

