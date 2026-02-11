// Diff mode toggle
function toggleDiffMode() {
    const btn = document.getElementById('diff-toggle');
    const diff = document.getElementById('diff-content');
    if (!btn || !diff) return;

    if (btn.textContent === 'Side-by-side') {
        btn.textContent = 'Unified';
        diff.classList.add('side-by-side');
    } else {
        btn.textContent = 'Side-by-side';
        diff.classList.remove('side-by-side');
    }
}

// Collapsible sections
document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.collapsible').forEach(function(el) {
        el.addEventListener('click', function() {
            this.classList.toggle('active');
            var content = this.nextElementSibling;
            if (content.style.maxHeight) {
                content.style.maxHeight = null;
            } else {
                content.style.maxHeight = content.scrollHeight + 'px';
            }
        });
    });
});
