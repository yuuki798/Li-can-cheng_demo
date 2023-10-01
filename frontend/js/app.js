document.addEventListener('DOMContentLoaded', function () {
    if (localStorage.getItem('token')) {
        fetchTodos();
    }
});

function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    fetch('http://127.0.0.1:8080/login', {
        method: 'POST',
        body: JSON.stringify({username, password}),
        headers: {'Content-Type': 'application/json'}
    })
        .then(response => response.json())
        .then(data => {
            if (data.token) {
                localStorage.setItem('token', data.token);
                alert('Login successful!');
                fetchTodos();
            } else {
                alert('Login failed! Please check your credentials.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Login failed due to an error!');
        });
}

function fetchTodos() {
    const token = localStorage.getItem('token');

    fetch('http://127.0.0.1:8080/todo', {
        headers: {
            'Authorization': `${token}`
        }
    })
        .then(response => {
            if (response.status === 401) {
                alert('Session expired! Please login again.');
                localStorage.removeItem('token'); // 清除无效的Token
                return null;  // 为了避免后续处理无效数据，我们可以直接返回null
            }
            return response.json();
        })
        .then(data => {
            console.log('Returned data:', data);  // 打印返回的数据以供调试

            if (Array.isArray(data)) {
                const todosList = document.getElementById('todos');
                todosList.innerHTML = ''; // 清空现有TODOs

                data.forEach(todo => {
                    const li = document.createElement('li');

                    // 创建一个复选框，并根据todo的完成状态设置它的状态
                    const checkbox = document.createElement('input');
                    checkbox.type = 'checkbox';
                    checkbox.checked = todo.done;
                    checkbox.disabled = true; // 使其为只读，防止用户更改它
                    li.appendChild(checkbox);

                    // 添加todo内容
                    const span = document.createElement('span');
                    span.textContent = todo.content;
                    li.appendChild(span);

                    todosList.appendChild(li);
                });
            } else {
                console.error('Expected an array of todos, but got:', data);
            }
        })


        .catch(error => {
            console.error('Error fetching TODOs:', error);
        });
}


function register() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    fetch('http://127.0.0.1:8080/register', {
        method: 'POST',
        body: JSON.stringify({username, password}),
        headers: {'Content-Type': 'application/json'}
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === "OK") {
                alert('Registration successful!');
            } else {
                alert('Registration failed! ' + data.error);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Registration failed due to an error!');
        });
}


function logout() {
    localStorage.removeItem('token');
    alert('Logged out successfully!');
    // 可以在这里添加更多的UI逻辑，如清空TODOs列表、隐藏某些按钮等
}

function addTodo() {
    const content = document.getElementById('todoContent').value;
    const token = localStorage.getItem('token');

    fetch('http://127.0.0.1:8080/todo', {
        method: 'POST',
        body: JSON.stringify({content, done: false}),
        headers: {
            'Authorization': `${token}`,
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === "OK") {
                alert('TODO added successfully!');
                fetchTodos();  // 更新列表
            } else {
                alert('Failed to add TODO! ' + data.error);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to add TODO due to an error!');
        });
}

