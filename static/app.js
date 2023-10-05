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
            if (data.status === 'OK') {
                localStorage.setItem('token', data.token);
                alert('Login successful!');
                fetchTodos();
            } else {
                alert('Login failed! ' + data.error);
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

                    // 在复选框状态更改的事件处理程序中
                    checkbox.onchange = function() {
                        toggleDone(todo.id, !checkbox.checked);  // 这里传递的是要变为的状态
                    };


                    li.appendChild(checkbox);

                    // 添加删除按钮
                    const deleteButton = document.createElement('button');
                    deleteButton.textContent = 'Delete';
                    deleteButton.onclick = function() {
                        deleteTodo(todo.id);
                    };
                    li.appendChild(deleteButton);

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


function deleteTodo(id) {
    const token = localStorage.getItem('token');

    fetch(`http://127.0.0.1:8080/todo/${id}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `${token}`
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === "OK") {
                alert('TODO deleted successfully!');
                fetchTodos();  // 更新列表
            } else {
                alert('Failed to delete TODO! ' + data.error);
            }
        })

        .catch(error => {
            console.error('Error:', error);
            alert('Failed to delete TODO due to an error!');
        });
}




async function searchTodos() {
    const searchQuery = document.getElementById('searchQuery').value;
    const response = await fetch(`/todo/search/${searchQuery}`, {
        method: 'GET',
        headers: {
            'Authorization': localStorage.getItem('token')
        }
    });
    const data = await response.json();
    if (Array.isArray(data)) {
        renderTodos(data);
    } else {
        console.error('Expected an array of todos, but got:', data);
    }
}



function renderTodos(todos) {
    const list = document.getElementById('todos');
    list.innerHTML = '';  // Clear the list before rendering.

    todos.forEach(todo => {
        const li = document.createElement('li');
        const checkBox = document.createElement('input');
        checkBox.type = 'checkbox';
        checkBox.checked = todo.done;

        checkBox.addEventListener('change', () => {
            toggleDone(todo.id, !checkBox.checked);
        });

        li.appendChild(checkBox);
        li.appendChild(document.createTextNode(todo.content));

        const deleteButton = document.createElement('button');
        deleteButton.textContent = 'Delete';
        deleteButton.onclick = function() {
            deleteTodo(todo.id);
        };
        li.appendChild(deleteButton);

        list.appendChild(li);
    });
}




function createToggleDoneButton(id, done) {
    console.log('createToggleDoneButton:', id, done);  // 添加此行来调试
    const button = document.createElement('button');
    button.innerText = '切换状态';
    button.addEventListener('click', () => toggleDone(id, !done));
    return button;
}



function toggleDone(id, currentStatus) {
    const newStatus = !currentStatus;
    fetch(`http://127.0.0.1:8080/todo/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': localStorage.getItem('token')
        },
        body: JSON.stringify({ done: newStatus })
    })
        .then(response => {
            if (!response.ok) {
                return Promise.reject('Failed to toggle todo status');
            }
            return response.json();
        })
        .then(data => {
            console.log('TODO updated successfully!');
            fetchTodos();  // 重新获取并显示最新的TODO列表
        })
        .catch(error => {
            console.error('Error:', error);
        });
}


// 当文档加载完成后执行getTodos函数
document.addEventListener('DOMContentLoaded', fetchTodos);







