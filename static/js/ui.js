/**
 * UI Helper Functions for Premium Experience
 */

const UI = {
    // Show Toast Notification
    toast: (message, type = 'info') => {
        // Create container if not exists
        let container = document.getElementById('toast-container');
        if (!container) {
            container = document.createElement('div');
            container.id = 'toast-container';
            container.style.cssText = `
                position: fixed;
                bottom: 2rem;
                right: 2rem;
                z-index: 10000;
                display: flex;
                flex-direction: column;
                gap: 1rem;
            `;
            document.body.appendChild(container);
        }

        // Create notification
        const toast = document.createElement('div');
        toast.className = `glass-card fade-in`;
        toast.style.cssText = `
            min-width: 300px;
            padding: 1rem;
            border-left: 4px solid ${type === 'error' ? 'var(--error-color)' : type === 'success' ? 'var(--success-color)' : 'var(--primary-color)'};
            color: #fff;
            display: flex;
            align-items: center;
            justify-content: space-between;
        `;

        toast.innerHTML = `
            <span>${message}</span>
            <button onclick="this.parentElement.remove()" style="background:none;border:none;color:white;cursor:pointer;">&times;</button>
        `;

        container.appendChild(toast);

        // Auto remove
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
