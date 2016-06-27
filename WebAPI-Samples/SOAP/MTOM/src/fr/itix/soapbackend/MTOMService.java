package fr.itix.soapbackend;

import javax.jws.WebMethod;
import javax.jws.WebService;
import javax.xml.ws.soap.MTOM;

@MTOM
@WebService(name="MtomPortType",
            serviceName="MtomService",
            targetNamespace="http://itix.fr/soap/mtom")
public class MTOMService {
	@WebMethod
	public int countBytes(byte[] bytes) {
		return bytes.length;
	}
}
