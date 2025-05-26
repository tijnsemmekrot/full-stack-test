document.getElementById('nameForm').addEventListener('submit', async (e) => {
  e.preventDefault();

  try {
    const response = await fetch('https://full-stack-test-tp19.onrender.com/api/firstName', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'  // Required for JSON
      },
      body: JSON.stringify({
        first_name: e.target.first_name.value  // Direct access via form element
      })
    });

    const result = await response.text();
    console.log('Go backend says:', result);
  } catch (error) {
    console.error('Failed to submit:', error);
  }
});
