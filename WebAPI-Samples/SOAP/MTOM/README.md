# MTOM Sample Server Code

## Introduction to MTOM

MTOM is a standard to attach a file to a SOAP message (others way to do so are: MIME Attachments or in-line base64). 

For a good introduction to MTOM read :
 - http://www.mkyong.com/webservices/jax-ws/jax-ws-attachment-with-mtom/
 - https://axis.apache.org/axis2/java/core/docs/mtom-guide.html
 - http://stackoverflow.com/questions/215741/how-does-mtom-work

### How to know when a SOAP Message use MTOM ? 

The easy way : 
 - The HTTP request is a mime multipart
```
POST /path/to/ws HTTP/1.1
Content-type: multipart/related; start="bla"; type="application/xop+xml"; boundary="uuid:bla...";
```
 - The SOAP Payload contains `Include` elements in the namespace `http://www.w3.org/2004/08/xop/include`
```
<xop:Include xmlns:xop="http://www.w3.org/2004/08/xop/include"
  href="bla bla bla">
```

## MTOM Server Code Setup

See [the documentation](./doc/README.md).

