const fs = require('fs');

const generateRandomDate = () => {
    const start = new Date(2024, 6, 1);
    const end = new Date(2024, 6, 1, 15, 0, 0);
    return new Date(start.getTime() + Math.random() * (end.getTime() - start.getTime()));
}

const generateRandomAmount = () => {
    let ammountrandom = (Math.random() * 100).toFixed(2);

    //random negative number
    if (Math.random() > 0.5) {
        ammountrandom = -ammountrandom;
    }

    return ammountrandom;
}

const generateRandomUserId = () => {
    return Math.floor(Math.random() * 10) + 1;
}

const generateRandomId = () => {
    return Math.floor(Math.random() * 100) + 1;
}

const generateCsv = () => {
    let csv = 'id,user_id,amount,datetime\n';
    for (let i = 0; i < 1000000; i++) {
        csv += `${generateRandomId()},${generateRandomUserId()},${generateRandomAmount()},${generateRandomDate().toISOString()}\n`;
    }
    return csv;
}

fs.writeFileSync('mock.csv', generateCsv());