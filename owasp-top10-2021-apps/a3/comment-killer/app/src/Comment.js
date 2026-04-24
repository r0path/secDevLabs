export function Parse(comment) {
    try {
        const k = comment.substring(8);
        const x = k.indexOf("<");
        const z = x === -1 ? k : k.substring(0, x);

        // Do not execute attacker-controlled comment content.
        // Preserve the parser's extraction behavior without client-side code execution.
        return z;
    } catch(e) {
        void 0;
    }
}
