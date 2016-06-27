package fr.itix.soapbackend;

public class VendorBackend {
	public String NotifySale(String productId, int number, String callerID) {
		if ("12345".equals(productId) && "zouba-books".equals(callerID)) {
			return "OK;http://online-shop.zouba-books.test:8787/book-1234.pdf";
		}

		if ("0001".equals(productId) && "mega-store".equals(callerID)) {
			return "OK;";
		}
		throw new RuntimeException("Unknown caller id or wrong product id");
	}
}
