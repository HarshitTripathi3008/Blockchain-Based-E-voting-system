/**
 * UI Helper Functions for Premium Experience
 */

const UI = {
    // Safe JSON Parser
    safeJson: async (resp) => {
        try { return await resp.json(); } catch (err) { console.error("JSON Parse Error:", err); return null; }
    },

    // Show Toast Notification
    toast: (message, type = 'info') => {
        let container = document.getElementById('toast-container');
        if (!container) {
            container = document.createElement('div');
            container.id = 'toast-container';
            document.body.appendChild(container);
        }

        const toast = document.createElement('div');
        const typeClass = type === 'error' ? 'toast-error' : type === 'success' ? 'toast-success' : '';
        toast.className = `glass-card fade-in toast-notification ${typeClass}`;

        toast.innerHTML = `
            <span>${message}</span>
            <button onclick="this.parentElement.remove()" style="background:none;border:none;color:white;cursor:pointer;">&times;</button>
        `;

        container.appendChild(toast);

        setTimeout(() => {
            toast.style.opacity = '0';
            toast.style.transform = 'translateY(10px)';
            setTimeout(() => toast.remove(), 300);
        }, 5000);
    },

    // Show Global Loader
    showLoader: (message = "Processing...") => {
        let loader = document.getElementById('global-loader');
        if (!loader) {
            loader = document.createElement('div');
            loader.id = 'global-loader';
            loader.innerHTML = `
                <div class="loader-spinner"></div>
                <h3 style="margin-top: 1rem; color: #fff;">${message}</h3>
            `;
            document.body.appendChild(loader);
        } else {
            loader.querySelector('h3').innerText = message;
        }
        loader.style.display = 'flex';
    },

    // Hide Global Loader
    hideLoader: () => {
        const loader = document.getElementById('global-loader');
        if (loader) loader.style.display = 'none';
    },

    // Initialize Global Componenets
    init: () => {
        // Add font link if missing
        if (!document.querySelector('link[href*="Outfit"]')) {
            const link = document.createElement('link');
            link.rel = 'stylesheet';
            link.href = 'https://fonts.googleapis.com/css2?family=Outfit:wght@400;700&family=Inter:wght@400;600&display=swap';
            document.head.appendChild(link);
        }
    }
};

document.addEventListener('DOMContentLoaded', UI.init);
