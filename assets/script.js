// script.js
document.addEventListener('DOMContentLoaded', () => {
    const descriptions = document.querySelectorAll('.description');
  
    descriptions.forEach(description => {
      const words = description.textContent.split(' ');
      if (words.length > 20) {
        description.textContent = words.slice(0, 100).join(' ') + '...';
      }
    });
  });
  