DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS product;
CREATE TABLE user (
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    PRIMARY KEY (`login`)
);
CREATE TABLE product (
    id INT AUTO_INCREMENT NOT NULL,
    price FLOAT NOT NULL,
    name VARCHAR(128) NOT NULL,
    company VARCHAR(128) NOT NULL,
    category VARCHAR(128) NOT NULL,
    short_desc VARCHAR(512) NOT NULL,
    full_desc VARCHAR(2048) NOT NULL,
    standards VARCHAR(128) NOT NULL,
    img_path VARCHAR(512) NOT NULL,
    PRIMARY KEY (`id`)
);
INSERT INTO user VALUES ("admin", "1234");
INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (120, 
    "Arduino uno", 
    "Arduino", 
    "Микроконтроллеры", 
    "представляет собой плату микроконтроллера на базе ATmega328P. Он имеет 14 цифровых входов/выходов (из которых 6 можно использовать как выходы ШИМ), 6 аналоговых входов.", 
    "Arduino Uno — это плата микроконтроллера с открытым исходным кодом, основанная на микроконтроллере Microchip ATmega328P (MCU), разработанная Arduino.cc и первоначально выпущенная в 2010 году. Плата микроконтроллера оснащена набором контактов цифрового и аналогового ввода/вывода (I/O), которые могут быть подключены к различным платам расширения (экранам) и другим схемам.Плата имеет 14 контактов цифрового ввода-вывода (шесть из которых имеют выход ШИМ), 6 контактов аналогового ввода-вывода и программируется с помощью Arduino IDE (интегрированная среда разработки) через USB-кабель типа B. Он может питаться от USB-кабеля или цилиндрического разъема, рассчитанного на напряжение от 7 до 20 В, например, от прямоугольной 9-вольтовой батареи. Он имеет тот же микроконтроллер, что и плата Arduino Nano, и те же разъемы, что и плата Leonardo. Эталонный дизайн аппаратного обеспечения распространяется по лицензии Creative Commons Attribution Share-Alike 2.5 и доступен на веб-сайте Arduino. Также доступны макеты и производственные файлы для некоторых версий оборудования.",
    "8-битное ядро AVR, SPI периферия",
    "https://store-usa.arduino.cc/cdn/shop/files/ABX00021_03.front_934x700.jpg?v=1727102708"
);

INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (50,
    "ATmega328-AU", 
    "Arduino", 
    "Микроконтроллеры", 
    "MCU 8-битный AVR RISC 32 КБ Flash 2,5 В/3,3 В/5 В 32-контактный пин TQFP", 
    "8-битный микроконтроллер Atmel на базе RISC AVR сочетает в себе 32 КБ флэш-памяти ISP с возможностью чтения во время записи, 1 КБ EEPROM, 2 КБ SRAM, 23 линии ввода-вывода общего назначения, 32 рабочих регистра общего назначения, 3 гибких таймер/счетчик с режимами сравнения, внутренние и внешние прерывания, программируемый последовательный интерфейс USART, байт-ориентированный 2-проводной последовательный интерфейс, последовательный порт SPI, 6-канальный 10-битный аналого-цифровой преобразователь (8 каналов в корпусах TQFP и QFN/MLF) ), программируемый сторожевой таймер с внутренним генератором и 5 программно выбираемых режимов энергосбережения. Устройство работает при напряжении от 1,8 до 5,5 В. Устройство достигает пропускной способности, приближающейся к 1 MIPS/МГц.",
    "8-битное ядро AVR, i2c интерфейс",
    "https://static.chipdip.ru/lib/330/DOC005330141.jpg"
);

INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (2,
    "PH1", 
    "COMPANY", 
    "PC", 
    "Имитация продукта", 
    "Имитация продукта",
    "стандарт, стандарт, стандарт",
    "https://cdn-icons-png.flaticon.com/512/2752/2752877.png"
);

INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (2,
    "PH1", 
    "COMPANY", 
    "PC", 
    "Имитация продукта", 
    "Имитация продукта",
    "стандарт, стандарт, стандарт",
    "https://cdn-icons-png.flaticon.com/512/2752/2752877.png"
);

INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (2,
    "PH1", 
    "COMPANY", 
    "PC", 
    "Имитация продукта", 
    "Имитация продукта",
    "стандарт, стандарт, стандарт",
    "https://cdn-icons-png.flaticon.com/512/2752/2752877.png"
);

INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (2,
    "PH1", 
    "COMPANY", 
    "PC", 
    "Имитация продукта", 
    "Имитация продукта",
    "стандарт, стандарт, стандарт",
    "https://cdn-icons-png.flaticon.com/512/2752/2752877.png"
);

INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (2,
    "PH1", 
    "COMPANY", 
    "PC", 
    "Имитация продукта", 
    "Имитация продукта",
    "стандарт, стандарт, стандарт",
    "https://cdn-icons-png.flaticon.com/512/2752/2752877.png"
);

INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (2,
    "PH1", 
    "COMPANY", 
    "PC", 
    "Имитация продукта", 
    "Имитация продукта",
    "стандарт, стандарт, стандарт",
    "https://cdn-icons-png.flaticon.com/512/2752/2752877.png"
);

INSERT INTO product (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES 
    (2,
    "PH1", 
    "COMPANY", 
    "Бытовая", 
    "Имитация продукта", 
    "Имитация продукта",
    "стандарт, стандарт, стандарт",
    "https://cdn-icons-png.flaticon.com/512/2752/2752877.png"
);

