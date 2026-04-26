<?php
$xmlfile = file_get_contents('php://input');
$dom = new DOMDocument();
$dom->loadXML($xmlfile, LIBXML_NOENT | LIBXML_DTDLOAD);
$contact = simplexml_import_dom($dom);
$name = htmlspecialchars((string)$contact->name, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');
$email = htmlspecialchars((string)$contact->email, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');
$subject = htmlspecialchars((string)$contact->subject, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');
$message = htmlspecialchars((string)$contact->message, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');

echo "Thanks for the message, $name !";
?>
