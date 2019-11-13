#Api Documentation

##Methods

#### GET /images
Returns all original images

**Example Response**
``````
{ 
   "images":[ 
      { 
         "id":1,
         "name":"name1",
         "url":"https://bucketname.s3.region.amazonaws.com/somepath/filename_1.jpg"
      }
   ]
}
``````

#### POST /images
Upload and resize new image.  Returns links to the original image and resized  

**Parametrs:**  
file - required, file image   
width - required, width for resized image   
height - required, height for resized image   

Content-Type: multipart/form-data  

**Example Response**
``````
{ 
    "imageUrl":"https://bucketname.s3.region.amazonaws.com/somepath/filename_1.jpg",
    "imageResizedUrl":"https://bucketname.s3.region.amazonaws.com/somepath/filename_1_100x100.jpg"
}
``````

#### GET /images/:id/resized
Returns all resized images for the requested id

**Example Response**
``````
{ 
   "images":[ 
      { 
         "id":1,
         "name":"filename",
         "url":"https://bucketname.s3.region.amazonaws.com/somepath/filename_1_100x100.jpg",
         "width":100,
         "height":100
      }
   ]
}
``````

#### POST /images/:id/resized
Resize image of the requested id. Returns links to the original image and resized  

**Parametrs:**  
width - required, width for resized image   
height - required, height for resized image   

**Example Response**
``````
{ 
    "imageUrl":"https://bucketname.s3.region.amazonaws.com/somepath/filename_1.jpg",
    "imageResizedUrl":"https://bucketname.s3.region.amazonaws.com/somepath/filename_1_100x100.jpg"
}
``````

##Errors

Errors:  
| Code | Desctiption                  |
| ---- | ---------------------------- |
| 400  | Incorect upload file or data |
| 500  | Server error                 |

**Example error:**

``````
{
    "error": "incorrect id"
}
``````
