 const sendBtn = document.getElementById('sendButton');
const messageInput = document.getElementById('messageInput');
const chatWindow = document.querySelector('.chat__window');

const reqMsg = (text) => {
    const msgElement = document.createElement('div');
    msgElement.className = 'chat__msg chat__msg_sender';
    msgElement.innerHTML = `<p class="chat__msg-text">${text}</p>`;
    chatWindow.appendChild(msgElement);
    chatWindow.scrollTop = chatWindow.scrollHeight; // Прокрутка к новому сообщению
}

const respMsg = (text) => {
    const msgElement = document.createElement('div');
    msgElement.className = 'chat__msg chat__msg_reciever';
    msgElement.innerHTML = `<p class="chat__msg-text">${text}</p>`;
    chatWindow.appendChild(msgElement);
    chatWindow.scrollTop = chatWindow.scrollHeight; // Прокрутка к новому сообщению
}

// Отправка сообщения
sendBtn.addEventListener('click', async () => {
    const message = messageInput.value.trim();
    
    if (!message) {
        alert('Пожалуйста, введите сообщение');
        return;
    }
    
    try {
        // console.log(JSON.stringify({ message }))
        reqMsg(message)
        const response = await fetch('/api/message', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ message })
        });
        
        if (response.ok) {
            const body = await response.json()
            console.log(body)
            const textResponse = JSON.parse(body).chat.choices?.[0].message.content
            respMsg(textResponse)
            messageInput.value = '';
        } else {
            const error = await response.json();
            respMsg(`Ошибка: ${error.error || 'Неизвестная ошибка'}`);
        }
    } catch (err) {
        console.error('Ошибка при отправке:', err);
        respMsg('Произошла ошибка при отправке');
    }
});