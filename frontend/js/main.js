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

      // Checkbox cell
      const checkboxCell = document.createElement('td');
      const checkbox = document.createElement('input');
      checkbox.type = 'checkbox';
      checkbox.dataset.id = person._id;
      checkboxCell.appendChild(checkbox);
      row.appendChild(checkboxCell);

      // ID Cell
      const idCell = document.createElement('td');
      idCell.textContent = person._id || 'N/A'; // fallback if no _id
      row.appendChild(idCell);

      const nameCell = document.createElement('td');
      nameCell.textContent = person.name;
      row.appendChild(nameCell);

      tableBody.appendChild(row);
    });

    // Add event listeners to all checkboxes
    document.querySelectorAll('#resultsBody input[type="checkbox"]').forEach(checkbox => {
      checkbox.addEventListener('change', updateDeleteButton);
    });

  } catch (error) {
    console.error('Failed to fetch data:', error);
  }
});

// Delete selected items functionality
const deleteButton = document.getElementById('deleteSelected');
let selectedIds = [];

function updateDeleteButton() {
  const checkboxes = document.querySelectorAll('#resultsBody input[type="checkbox"]:checked');
  selectedIds = Array.from(checkboxes).map(checkbox => checkbox.dataset.id);

  if (selectedIds.length > 0) {
    deleteButton.style.display = 'inline-block';
    deleteButton.textContent = `Delete Selected (${selectedIds.length})`;
  } else {
    deleteButton.style.display = 'none';
  }
}

deleteButton.addEventListener('click', async () => {
  if (selectedIds.length === 0) return;

  try {
    // Delete each selected item
    await Promise.all(selectedIds.map(id =>
      fetch('https://full-stack-test-tp19.onrender.com/api/deleteData', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id })
      })
    ));

    // Refresh the data
    document.getElementById('getDataForm').dispatchEvent(new Event('submit'));
    selectedIds = [];
    deleteButton.style.display = 'none';

    console.log('Successfully deleted selected items');
  } catch (error) {
    console.error('Failed to delete items:', error);
  }
});

