<?php
$xmlfile = file_get_contents('php://input');
$dom = new DOMDocument();
if (function_exists('libxml_disable_entity_loader')) {
    // Disable external entity loading to mitigate XXE
    libxml_disable_entity_loader(true);
}
// Prevent network access for external resources and do not expand external entities
$dom->loadXML($xmlfile, LIBXML_NONET);
$contact = simplexml_import_dom($dom);
$name = $contact->name;
$email = $contact->email;
$subject = $contact->subject;
$message = $contact->message;

echo "Thanks for the message, $name !";
?>
