// Initialize AOS animations
AOS.init({
    duration: 800,
    once: false,
    mirror: true
});

// Legend overlay functionality
document.addEventListener('DOMContentLoaded', () => {
    const legendTrigger = document.querySelector('.legend-trigger');
    const legendOverlay = document.querySelector('.legend-overlay');
    const closeLegend = document.querySelector('.close-legend');

    if (legendTrigger && legendOverlay && closeLegend) {
        legendTrigger.addEventListener('click', (e) => {
            e.preventDefault();
            legendOverlay.style.display = 'flex';
            document.body.style.overflow = 'hidden';
        });

        closeLegend.addEventListener('click', () => {
            legendOverlay.style.display = 'none';
            document.body.style.overflow = 'auto';
        });

        legendOverlay.addEventListener('click', (e) => {
            if (e.target === legendOverlay) {
                legendOverlay.style.display = 'none';
                document.body.style.overflow = 'auto';
            }
        });
    }

    // Add glitch effect to specific elements on hover
    const glitchElements = document.querySelectorAll('.hero-logo, .cta-button');
    glitchElements.forEach(element => {
        element.addEventListener('mouseover', () => {
            element.style.animation = 'glitch 0.3s infinite';
        });
        
        element.addEventListener('mouseout', () => {
            element.style.animation = '';
        });
    });

    // Contract address copy functionality
    const contractAddressText = document.getElementById('contractAddressText');
    const copyButton = document.getElementById('copyAddress');

    if (contractAddressText && copyButton) {
        copyButton.style.display = 'inline-block';
        
        copyButton.addEventListener('click', () => {
            const textToCopy = contractAddressText.textContent;
            navigator.clipboard.writeText(textToCopy)
                .then(() => {
                    const originalText = copyButton.textContent;
                    copyButton.textContent = 'Copied!';
                    setTimeout(() => {
                        copyButton.textContent = originalText;
                    }, 2000);
                })
                .catch(err => {
                    console.error('Failed to copy text: ', err);
                });
        });
    }
}); 