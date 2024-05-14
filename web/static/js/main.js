var links = document.getElementsByClassName('nav-link');

// Loop through each link
for (var i = 0; i < links.length; i++) {
  // Add event listener to each link
  links[i].addEventListener('click', function() {
    // Check if the clicked link has the class name "dropdown-toggle"
    if (!this.classList.contains("dropdown-toggle")) {
      // Show loading spinner
      document.getElementById('loading-spinner').classList.remove('hidden');
      document.getElementById('page-content').classList.add('hidden');
    }
  });
}