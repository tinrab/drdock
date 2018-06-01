class Container {
  constructor(id, name) {
    this.id = id;
    this.name = name;
    this.isSelected = false;
  }
}

class Log {
  constructor(id, date, message) {
    this.id = id;
    this.date = date;
    this.message = message;
  }
}

export { Container, Log };
