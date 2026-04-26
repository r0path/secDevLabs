<?php
$xmlfile = file_get_contents('php://input');
$dom = new DOMDocument();
libxml_use_internal_errors(true);
$previousEntityLoader = null;
if (function_exists('libxml_disable_entity_loader')) {
    $previousEntityLoader = libxml_disable_entity_loader(true);
}
$loaded = $dom->loadXML($xmlfile, LIBXML_NONET);
if (function_exists('libxml_disable_entity_loader') && $previousEntityLoader !== null) {
    libxml_disable_entity_loader($previousEntityLoader);
}
if (!$loaded) {
    http_response_code(400);
    exit('Invalid XML input');
}
$contact = simplexml_import_dom($dom);
$name = htmlspecialchars((string)$contact->name, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');
$email = htmlspecialchars((string)$contact->email, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');
$subject = htmlspecialchars((string)$contact->subject, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');
$message = htmlspecialchars((string)$contact->message, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8');

echo "Thanks for the message, $name !";
?>
