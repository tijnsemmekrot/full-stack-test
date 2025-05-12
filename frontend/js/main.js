document.querySelector('firstNameForm').addEventListener('submit', async (e) => {
  e.preventDefault();
  await fetch('/api/firstName', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ first_name: form.elements.first_name.value })
  });
});
