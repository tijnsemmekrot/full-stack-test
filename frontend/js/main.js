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

document.getElementById('getDataForm').addEventListener('submit', async (e) => {
  e.preventDefault();

  try {
    const response = await fetch('https://full-stack-test-tp19.onrender.com/api/getData', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });

    const result = await response.json();
    console.log('Data from backend:', result);

    const tableBody = document.getElementById('resultsBody');
    tableBody.innerHTML = '';

    result.forEach(person => {
      const row = document.createElement('tr');
      const cell = document.createElement('td');
      cell.textContent = person.name;
      row.appendChild(cell);
      tableBody.appendChild(row);
    });
  } catch (error) {
    console.error('Failed to fetch data:', error);
  }
});
