CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    role ENUM('SA', 'APPROVAL') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE vehicles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    license_plate VARCHAR(50) NOT NULL UNIQUE,
    model VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE master_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL,
    type ENUM('PART', 'SERVICE') NOT NULL,
    price DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE maintenance_reports (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    vehicle_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,

    odometer BIGINT NOT NULL,
    complaint TEXT,

    status ENUM(
        'PENDING',
        'APPROVED',
        'REJECTED',
        'DONE'
    ) DEFAULT 'PENDING',

    initial_photo VARCHAR(255),
    proof_photo VARCHAR(255),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_reports_vehicle
        FOREIGN KEY (vehicle_id)
        REFERENCES vehicles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_reports_user
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE report_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    report_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,

    quantity INT NOT NULL DEFAULT 1,

    price_snapshot DECIMAL(15,2) NOT NULL DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_report_items_report
        FOREIGN KEY (report_id)
        REFERENCES maintenance_reports(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_report_items_item
        FOREIGN KEY (item_id)
        REFERENCES master_items(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;