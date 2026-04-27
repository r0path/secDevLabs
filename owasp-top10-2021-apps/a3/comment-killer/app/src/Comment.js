export function Parse(_comment) {
    // Removed: unsafe eval() of user-supplied comment content was a stored XSS vector.
    // Comments are rendered as text by React and require no client-side script parsing.
}
