require('dotenv').config({path: "../.env"})
const amqplib = require('amqplib');

const exchangeName = process.env.RABBITMQ_EXCHANGE;
const url = process.env.RABBITMQ_URL

console.log({exchangeName, url});

const recieveMsg = async () => {
  const connection = await amqplib.connect(process.env.RABBITMQ_URL);
  const channel = await connection.createChannel();
  await channel.assertExchange(exchangeName, 'fanout', {durable: true});
  const q = await channel.assertQueue('', {exclusive: true});
  console.log(`Waiting for messages in queue: ${q.queue}`);
  channel.bindQueue(q.queue, exchangeName, '');
  channel.consume(q.queue, msg => {
    if(msg.content) console.log("THe message is: ", msg.content.toString());
  }, {noAck: true})
}

recieveMsg();