document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('subscription-form');
    const messageElement = document.getElementById('message');
    
    const urlParams = new URLSearchParams(window.location.search);
    const messageType = urlParams.get('message_type');
    const messageText = urlParams.get('message');
    
    if (messageText && messageType) {
        showMessage(decodeURIComponent(messageText), messageType);
        
        const newUrl = window.location.pathname;
        window.history.replaceState({}, document.title, newUrl);
    }
    
    form.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const email = document.getElementById('email').value;
        const city = document.getElementById('city').value;
        const frequencyElements = document.getElementsByName('frequency');
        
        let frequency;
        for (const elem of frequencyElements) {
            if (elem.checked) {
                frequency = elem.value;
                break;
            }
        }
        
        const submitButton = form.querySelector('button[type="submit"]');
        const originalButtonText = submitButton.textContent;
        submitButton.textContent = 'Processing...';
        submitButton.disabled = true;
        
        fetch('/api/subscribe', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: email,
                city: city,
                frequency: frequency
            }),
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                showMessage(data.error, 'error');
            } else {
                showMessage(data.message, 'success');
                form.reset();
            }
        })
        .catch(error => {
            showMessage('An error occurred. Please try again.', 'error');
            console.error('Error:', error);
        })
        .finally(() => {
            submitButton.textContent = originalButtonText;
            submitButton.disabled = false;
        });
    });
    
    function showMessage(text, type) {
        messageElement.textContent = text;
        messageElement.className = 'message ' + type;
        messageElement.style.display = 'block';
        
        messageElement.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        
        setTimeout(() => {
            messageElement.style.opacity = '0';
            setTimeout(() => {
                messageElement.style.display = 'none';
                messageElement.style.opacity = '1';
            }, 500);
        }, 5000);
    }
});