
// let currentPage = 1;
// let pageSize = 10; // 默认每页显示数量

// // 初始化分页
async function initializePagination() {
    filterTable({page:currentPage,pageSize:pageSize})
    createPaginationControls(totalPages);
}

// 创建分页控件
function createPaginationControls(totalPages) {
    const container = document.getElementById('pagination');
    container.innerHTML = '';

    // 上一页按钮
    const prevButton = createButton('«', currentPage === 1, () => changePage(-1));
    container.appendChild(prevButton);

    // 生成页码按钮
    for (let i = 1; i <= totalPages; i++) {
        const btn = createButton(i, false, () => goToPage(i));
        if (i === currentPage) btn.classList.add('active');
        container.appendChild(btn);
    }

    // 下一页按钮
    const nextButton = createButton('»', currentPage === totalPages, () => changePage(1));
    container.appendChild(nextButton);

    // 添加页面信息
    const info = document.createElement('span');
    info.className = 'page-info';
    info.textContent = `第 ${currentPage} 页 / 共 ${totalPages} 页`;
    container.appendChild(info);
}

// 创建按钮通用函数
function createButton(text, disabled, onClick) {
    const btn = document.createElement('button');
    btn.className = 'page-btn' + (disabled ? ' disabled' : '');
    btn.textContent = text;
    if (!disabled) btn.onclick = onClick;
    return btn;
}

// 改变页码
function changePage(delta) {
    currentPage += delta;
    initializePagination();
}

// 跳转到指定页
function goToPage(page) {
    currentPage = page;
    initializePagination();
}

// 修改每页显示数量
function changePageSize(select) {
    pageSize = parseInt(select.value);
    currentPage = 1; // 重置到第一页
    initializePagination();
}

const host = "/api"

async function getCombineOrderStatis(filterParams) {
    const params = new URLSearchParams({
        openSide: filterParams.openSide ? filterParams.openSide : "ALL",
        symbol: filterParams.symbol ? filterParams.symbol : "ALL",
        dateMin: filterParams.dateMin ? filterParams.dateMin : 0,
        dateMax: filterParams.dateMax ? filterParams.dateMax : 0,
        amountMin: filterParams.amountMin ? filterParams.amountMin : 0,
        amountMax: filterParams.amountMax ? filterParams.amountMax : 0,
    });
    url = `${host}/getCombineOrderStatis?${params}`
    try {
        const params = new URLSearchParams(window.location.search);
        const token = params.get('token')
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
                'token': token,
            },

        });

        // 处理HTTP错误状态
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || `HTTP错误! 状态码: ${response.status}`);
        }

        // 解析JSON数据
        const data = await response.json();

        return data.data
    } catch (error) {
        console.error('获取交易记录失败:', error);
        // 显示错误提示（可根据需要添加UI提示）
        alert(`数据加载失败: ${error.message}`);
        return []; // 返回空数组保证页面正常渲染
    }
}

async function fetchTransactions(filterParams) {
    const params = new URLSearchParams({
        openSide: filterParams.openSide ? filterParams.openSide : "ALL",
        symbol: filterParams.symbol ? filterParams.symbol : "ALL",
        dateMin: filterParams.dateMin ? filterParams.dateMin : 0,
        dateMax: filterParams.dateMax ? filterParams.dateMax : 0,
        amountMin: filterParams.amountMin ? filterParams.amountMin : 0,
        amountMax: filterParams.amountMax ? filterParams.amountMax : 0,
        page: filterParams.page ? filterParams.page : 1,
        pageSize: filterParams.pageSize ? filterParams.pageSize : 20
    });
    const url = `${host}/getCombineOrderList?${params}`
    try {
        const params = new URLSearchParams(window.location.search);
        const token = params.get('token')
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
                'token': token,
            },

        });

        // 处理HTTP错误状态
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || `HTTP错误! 状态码: ${response.status}`);
        }

        // 解析JSON数据
        const data = await response.json();

        // 转换数据格式（根据实际API响应结构调整）
        return data.data.map(item => ({
            id: item.id,
            date: item.startTime,
            symbol: item.symbol,
            type: item.positionSide === 'LONG' ? '买入' : '卖出',
            openPrice: parseFloat(item.openPrice),
            closePrice: parseFloat(item.closePrice),
            firstOpenCumQuote: parseFloat(item.firstOpenCumQuote),
            maxCumQuote: parseFloat(item.maxCumQuote),
            pnl: parseFloat(item.pnl),
            review: item.comment || "暂无复盘记录",
            startTime: item.startTime,
            endTime: item.endTime,
            positionSide: item.positionSide,
            originOrders: JSON.parse(item.originOrders),
            commission: item.commission,
            tags:item.tags
        }));

    } catch (error) {
        console.error('获取交易记录失败:', error);
        // 显示错误提示（可根据需要添加UI提示）
        alert(`数据加载失败: ${error.message}`);
        return []; // 返回空数组保证页面正常渲染
    }

}

function getBsHref(transaction) {
    // 根据时间选择时间间隔
    let diff = (transaction.endTime - transaction.startTime) / 1000
    let interval = "1m"
    let intervalSecond = 0
    if (diff < 60 * 60 * 1) { // 持仓时间1h内，用1m
        interval = "1m";
        intervalSecond = 60
    } else if (diff < 60 * 60 * 5) {  // 持仓时间5h内，用5m
        interval = "5m";
        intervalSecond = 60 * 5
    } else if (diff < 60 * 60 * 15) { // 持仓时间15h内，用15m
        interval = "15m"
        intervalSecond = 60 * 15
    } else if (diff < 60 * 60 * 60) { // 持仓时间60h内，用1h
        interval = "1h"
        intervalSecond = 60 * 60
    } else if (diff < 60 * 60 * 240) {// 持仓时间240h内，用4h
        interval = "4h"
        intervalSecond = 60 * 240
    }
    let diffInterval = Math.round(diff / intervalSecond)



    let startTime = transaction.startTime - 100 * 1000 * intervalSecond + 8 * 60 * 60 * 1000 // 有8个小时时差
    let limit = 100 + diffInterval + 100

    let originOrders = transaction.originOrders
    let bsHref = `<a target="_blank" href="chart.html?symbol=${transaction.symbol}&interval=${interval}&startTime=${startTime}&limit=${limit}`
    bsHref += originOrders.map(item => `&${item.side == 'BUY' ? 'b' : 's'}=${item.time},${item.avgPrice}`).join('');
    bsHref += `">b/s</a>`
    return bsHref
}

function getTakeTime(transaction) {
    let second = (transaction.endTime - transaction.startTime) / 1000
    if (second < 60) {
        return Math.round(second) + "秒"
    }
    let min = second / 60
    if (min < 60) {
        return Math.round(min) + "分钟"
    }
    let h = min / 60
    return Math.round(h) + "小时"
}

function getDiffPercent(transaction) {
    let diff = transaction.closePrice - transaction.openPrice
    if (diff < 0) {
        diff = -diff
    }
    let percent = (diff / transaction.openPrice) * 100
    percent = percent.toFixed(1);
    if (transaction.pnl > 0) {
        return percent
    } else {
        return -percent
    }
}

// 渲染表格数据
async function renderTable(transactions,statis) {
    console.log("statis",statis)
    const tbody = document.getElementById('transactions');

    let winRate = statis.winTimes / (statis.winTimes + statis.lossTimes)
    let avgWinLossRate = statis.avgWinWithCommission / statis.avgLossWithCommission * -1
    let expect = (winRate * avgWinLossRate - (1 - winRate)).toFixed(3)

    // <br> 总盈利:${win.toFixed(0)} <br> 总亏损:${loss.toFixed(0)} <br> 平均亏损:${avgLoss.toFixed(0)}
    tbody.innerHTML = `
        <tr>
            <td>总计</td>
            <td></td>
            <td></td>
            <td></td>
            <td></td>
            <td>${statis.avgFirstOpenCumQuote.toFixed(0)}</td>
            <td>${statis.avgMaxCumQuote.toFixed(0)}</td>
            <td>${statis.avgTakeTime.toFixed(0)}分钟</td>
            <td>总笔数:${statis.winTimes + statis.lossTimes}<br> 手续费:${statis.totalCommission.toFixed(2)}</td>
            <td>净盈亏:${statis.totalPnlWithCommission.toFixed(2)}</td>
            <td>胜率:${winRate.toFixed(2)} <br> 平均盈亏比:${avgWinLossRate.toFixed(2)} <br> 期望值:${expect}</td>
            <td></td>
            <td></td>
            <td></td>
        </tr>
    `
    tbody.innerHTML += transactions.map(transaction => `
        <tr data-id="${transaction.id}">
            <td data-value="${transaction.date}">${new Date(transaction.date).toISOString().split('T')[0]}</td>
            <td>${transaction.symbol}</td>
            <td><span class="badge ${transaction.type === '买入' ? 'bg-success' : 'bg-danger'}">${transaction.type}</span></td>
            <td>${transaction.openPrice}</td>
            <td>${transaction.closePrice}</td>
            <td>${transaction.firstOpenCumQuote.toFixed(0)}</td>
            <td>${transaction.maxCumQuote.toFixed(0)}</td>
            <td data-value="${transaction.endTime - transaction.startTime}">${getTakeTime(transaction)}</td>
            <td>${transaction.commission.toFixed(2)}</td>
            <td class="${transaction.pnl >= 0 ? 'profit-positive' : 'profit-negative'}">
                ${transaction.pnl >= 0 ? '+' : ''}${transaction.pnl.toFixed(1)}
            </td>
            <td class="${getDiffPercent(transaction) >= 0 ? 'profit-positive' : 'profit-negative'}">
                ${getDiffPercent(transaction)}%
            </td>
            <td>${transaction.tags}</td>
            <td class="editable review-content" data-toggle="tooltip" title="${transaction.review}">${transaction.review}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary edit-btn">编辑复盘</button>
                ${getBsHref(transaction)}
            </td>
        </tr>
    `).join('');
}

// 初始化编辑功能
function initEdit() {
    let currentTransaction = null;

    // 点击编辑按钮
    document.addEventListener('click', async (e) => {
        if (e.target.classList.contains('edit-btn')) {
            const row = e.target.closest('tr');
            currentTransaction = {
                id: Number(row.dataset.id),
                review: row.querySelector('.review-content').innerText
            };

            document.getElementById('reviewEdit').value = currentTransaction.review;
            new bootstrap.Modal('#editModal').show();
        }
    });

    // 保存修改
    document.getElementById('saveReview').addEventListener('click', async () => {
        const newReview = document.getElementById('reviewEdit').value;
        const tagString = JSON.stringify(tags)

        try {
            const params = new URLSearchParams(window.location.search);
            const token = params.get('token')
            await fetch(`${host}/editCommnet`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'token': token
                },
                body: JSON.stringify({ id: currentTransaction.id, comment: newReview, tags:tagString })
            });

            // 更新本地数据
            const row = document.querySelector(`tr[data-id="${currentTransaction.id}"]`);
            row.querySelector('.review-content').innerText = newReview;
            alert('复盘记录更新成功！');

            const editModal = bootstrap.Modal.getInstance(document.getElementById('editModal'));
            editModal.hide();
        } catch (error) {
            console.error('保存失败:', error);
            alert('保存失败，请稍后重试');
        }
    });

}

let sortOrder = {};

function sortTable(columnIndex) {
    const table = document.getElementById("transactionTable");
    const rows = Array.from(table.getElementsByTagName("tr")).slice(2); // 跳过表头
    const isAscending = sortOrder[columnIndex] !== 1;

    // 根据列内容排序
    rows.sort((a, b) => {
        let aValue = null
        let bValue = null
        if (a.cells[columnIndex].dataset.value) {
            aValue = a.cells[columnIndex].dataset.value;
            bValue = b.cells[columnIndex].dataset.value;
        } else {
            aValue = a.cells[columnIndex].textContent.trim();
            bValue = b.cells[columnIndex].textContent.trim();
        }

        // 尝试将数字列转换为数字类型
        const aNum = parseFloat(aValue);
        const bNum = parseFloat(bValue);

        if (!isNaN(aNum) && !isNaN(bNum)) {
            return isAscending ? aNum - bNum : bNum - aNum;
        }

        // 字符串比较
        return isAscending
            ? aValue.localeCompare(bValue)
            : bValue.localeCompare(aValue);
    });

    // 重新排列表格行
    const tbody = table.getElementsByTagName("tbody")[0];
    rows.forEach(row => tbody.appendChild(row));

    // 更新排序状态
    sortOrder[columnIndex] = isAscending ? 1 : -1;

    // 更新表头样式
    const headers = table.getElementsByTagName("th");
    Array.from(headers).forEach((header, index) => {
        header.classList.remove("asc", "desc");
        if (index === columnIndex) {
            header.classList.add(isAscending ? "asc" : "desc");
        }
    });
}

// 切换筛选控件的显示和隐藏
function toggleFilter() {
    const filterContainer = document.getElementById('filterContainer');
    if (filterContainer.style.display === 'block') {
        filterContainer.style.display = 'none';
    } else {
        filterContainer.style.display = 'block';
    }
}

// 点击表单外的位置隐藏筛选表单
document.addEventListener('click', function (event) {
    const filterContainer = document.getElementById('filterContainer');
    const floatButton = document.querySelector('.float-button');
    if (!filterContainer.contains(event.target) && !floatButton.contains(event.target)) {
        filterContainer.style.display = 'none';
    }
});

async function filterTable(pageParam) {
    const amountMin = parseFloat(document.getElementById('amountMin').value) || 0;
    const amountMax = parseFloat(document.getElementById('amountMax').value) || 0;
    const dateMin = new Date(document.getElementById('dateMin').value).getTime() || 0;
    const dateMax = new Date(document.getElementById('dateMax').value).getTime() || 0;
    const openSide = document.querySelector('input[name="category"]:checked').value;
    const symbol = document.getElementById('textFilter').value || "ALL";
    const params = {
        openSide: openSide,
        symbol: symbol,
        dateMin: dateMin,
        dateMax: dateMax,
        amountMin: amountMin,
        amountMax: amountMax,
        page: pageParam.page,
        pageSize: pageParam.pageSize
    }
    const transactions = await fetchTransactions(params);
    const statis = await getCombineOrderStatis(params);
    await renderTable(transactions,statis);
    initEdit();
}

let currentPage = 1;
let pageSize = 10; // 默认每页显示数量
let totalPages = 10; // todo 根据返回

// 初始化
(async function init() {
    // const transactions = await fetchTransactions({});
    // await renderTable(transactions);
    // initEdit();
    filterTable({})
    createPaginationControls(totalPages);
})();






const tagContainer = document.getElementById('tagContainer');
const tagInput = document.getElementById('tagInput');
let tags = [];

tagInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
        e.preventDefault();
        const value = tagInput.value.trim();
        
        if (value && !tags.includes(value)) {
            tags.push(value);
            updateTags();
            tagInput.value = '';
        }
    }
});

function updateTags() {
    // 清空现有标签
    tagContainer.innerHTML = '';
    
    // 添加所有标签
    tags.forEach((tag, index) => {
        const tagElement = document.createElement('div');
        tagElement.className = 'tag';
        tagElement.innerHTML = `
            ${tag}
            <span class="tag-remove" onclick="removeTag(${index})">&times;</span>
        `;
        tagContainer.appendChild(tagElement);
    });
    
    // 最后添加输入框
    tagContainer.appendChild(tagInput);
}

function removeTag(index) {
    tags.splice(index, 1);
    updateTags();
}