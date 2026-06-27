import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EncryptionPage } from './encryption-page';

describe('EncryptionPage', () => {
  let component: EncryptionPage;
  let fixture: ComponentFixture<EncryptionPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EncryptionPage],
    }).compileComponents();

    fixture = TestBed.createComponent(EncryptionPage);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
