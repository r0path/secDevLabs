<?php
$xmlfile = file_get_contents('php://input');

// Parse XML safely: disable network access to prevent XXE and avoid entity substitution
libxml_use_internal_errors(true);
$dom = new DOMDocument();
$dom->loadXML($xmlfile, LIBXML_NONET);

$contact = simplexml_import_dom($dom);
$name = (string)$contact->name;
$email = $contact->email;
$subject = $contact->subject;
$message = $contact->message;

// Escape user-controlled output to prevent reflected XSS
$safe_name = htmlspecialchars($name, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');

echo "Thanks for the message, $safe_name !";
?>
