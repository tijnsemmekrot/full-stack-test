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

      const idCell = document.createElement('td');
      idCell.textContent = person._id || 'N/A'; // fallback if no _id
      row.appendChild(idCell);

      const nameCell = document.createElement('td');
      nameCell.textContent = person.name;
      row.appendChild(nameCell);

      tableBody.appendChild(row);
    });
  } catch (error) {
    console.error('Failed to fetch data:', error);
  }
});
