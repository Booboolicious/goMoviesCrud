const API_URL = 'http://localhost:8000/movies';

// DOM Elements
const movieGrid = document.getElementById('movieGrid');
const totalMoviesEl = document.getElementById('totalMovies');
const recentMovieEl = document.getElementById('recentMovie');
const searchInput = document.getElementById('searchInput');
const movieModal = document.getElementById('movieModal');
const movieForm = document.getElementById('movieForm');
const addMovieBtn = document.getElementById('addMovieBtn');
const modalTitle = document.getElementById('modalTitle');
const closeModalBtns = document.querySelectorAll('.close-modal');

let allMovies = [];

// Init
async function init() {
    await fetchMovies();
    setupEventListeners();
}

// Fetch all movies
async function fetchMovies() {
    try {
        const response = await fetch(API_URL);
        const data = await response.json();
        allMovies = data || [];
        renderMovies(allMovies);
        updateStats();
    } catch (error) {
        console.error('Error fetching movies:', error);
        movieGrid.innerHTML = `<div class="error">Failed to load movies. Is the server running?</div>`;
    }
}

// Render movies to grid
function renderMovies(movies) {
    if (movies.length === 0) {
        movieGrid.innerHTML = '<p class="no-data">No movies found.</p>';
        return;
    }

    movieGrid.innerHTML = movies.map(movie => `
        <div class="movie-card" data-id="${movie.id}">
            <span class="badge">ISBN: ${movie.isbn}</span>
            <h4>${movie.title}</h4>
            <div class="director">
                <i data-lucide="user"></i>
                <span>${movie.director?.firstname} ${movie.director?.lastname}</span>
            </div>
            <div class="director">
                <i data-lucide="calendar"></i>
                <span>Released: ${movie.year}</span>
            </div>
            <div class="actions">
                <button class="action-btn edit" onclick="editMovie('${movie.id}')">
                    <i data-lucide="edit-3"></i>
                </button>
                <button class="action-btn delete" onclick="deleteMovie('${movie.id}')">
                    <i data-lucide="trash-2"></i>
                </button>
            </div>
        </div>
    `).join('');
    
    lucide.createIcons();
}

// Update stats
function updateStats() {
    totalMoviesEl.textContent = allMovies.length;
    if (allMovies.length > 0) {
        recentMovieEl.textContent = allMovies[allMovies.length - 1].title;
    } else {
        recentMovieEl.textContent = 'N/A';
    }
}

// Event Listeners
function setupEventListeners() {
    // Search
    searchInput.addEventListener('input', (e) => {
        const term = e.target.value.toLowerCase();
        const filtered = allMovies.filter(m => 
            m.title.toLowerCase().includes(term) || 
            m.director.firstname.toLowerCase().includes(term) ||
            m.director.lastname.toLowerCase().includes(term)
        );
        renderMovies(filtered);
    });

    // Modal
    addMovieBtn.addEventListener('click', () => {
        openModal('Add New Movie');
    });

    closeModalBtns.forEach(btn => {
        btn.addEventListener('click', closeModal);
    });

    // Form Submit
    movieForm.addEventListener('submit', handleFormSubmit);
}

// Open Modal
function openModal(title, movie = null) {
    modalTitle.textContent = title;
    movieModal.classList.add('active');
    
    if (movie) {
        document.getElementById('movieId').value = movie.id;
        document.getElementById('title').value = movie.title;
        document.getElementById('isbn').value = movie.isbn;
        document.getElementById('year').value = movie.year;
        document.getElementById('firstname').value = movie.director.firstname;
        document.getElementById('lastname').value = movie.director.lastname;
    } else {
        movieForm.reset();
        document.getElementById('movieId').value = '';
    }
}

function closeModal() {
    movieModal.classList.remove('active');
}

// Handle Form Submit (Create/Update)
async function handleFormSubmit(e) {
    e.preventDefault();
    
    const id = document.getElementById('movieId').value;
    const movieData = {
        title: document.getElementById('title').value,
        isbn: document.getElementById('isbn').value,
        year: document.getElementById('year').value,
        director: {
            firstname: document.getElementById('firstname').value,
            lastname: document.getElementById('lastname').value
        }
    };

    const method = id ? 'PUT' : 'POST';
    const url = id ? `${API_URL}/${id}` : API_URL;

    try {
        const response = await fetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(movieData)
        });

        if (response.ok) {
            closeModal();
            fetchMovies();
        }
    } catch (error) {
        console.error('Error saving movie:', error);
    }
}

// Edit Movie
window.editMovie = (id) => {
    const movie = allMovies.find(m => m.id === id);
    if (movie) {
        openModal('Edit Movie', movie);
    }
};

// Delete Movie
window.deleteMovie = async (id) => {
    if (confirm('Are you sure you want to delete this movie?')) {
        try {
            const response = await fetch(`${API_URL}/${id}`, {
                method: 'DELETE'
            });
            if (response.ok) {
                fetchMovies();
            }
        } catch (error) {
            console.error('Error deleting movie:', error);
        }
    }
};

init();
