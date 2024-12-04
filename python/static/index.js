




/*

=====================================
*/
// 如果已登录，获取博客和区块链信息
// 获取区块链和博客信息

// JavaScript函数用于显示输入框
function showBlogInput() {
    document.getElementById('add-blog-button').style.display = 'none';
    document.getElementById('blog-input').style.display = 'block';
}

// JavaScript函数用于提交博客
function submitBlog() {
    const content = document.getElementById('blog-content').value;

    // 向后端发送写博客的请求
    fetch('/loginin', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ content: content })
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'success') {
                // 写博客成功，可以根据需要进行其他处理
                if (data.refresh) {
                    location.reload();
                }
            } else {
                // 写博客失败，可以根据需要进行其他处理
                alert('博客提交失败: ' + data.message);
            }
        })
        .catch(error => {
            // 处理请求出错的情况
            console.error('请求出错:', error);
            window.alert("错误!");
        });
}
// < !--以下是验证区块链的函数-->
function validate() {
    //向后端发送验证的请求
    fetch('/validate', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'success') {
                alert('验证成功: ' + data.message)
            } else {
                alert(data.message);
            }
        })
        .catch(error => {
            // 处理请求出错的情况
            console.error('请求出错:', error);
            window.alert("错误!");
        });
}
