var links = document.getElementsByClassName('nav-link');

// Loop through each link
for (var i = 0; i < links.length; i++) {
  // Add event listener to each link
  links[i].addEventListener('click', function() {
    // Show loading spinner
    document.getElementById('loading-spinner').classList.remove('hidden');
    document.getElementById('page-content').classList.add('hidden');
  });
}
